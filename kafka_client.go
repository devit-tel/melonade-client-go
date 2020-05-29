package melonade_client_go

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/devit-tel/goerror"
)

type KafkaClient struct {
	namespace    string
	kafkaServers []string
	config       *sarama.Config
	producer     sarama.AsyncProducer
}

type Watcher struct {
	ChanTransaction    chan *EventTransaction
	ChanTransactionErr chan *EventTransactionError
	ChanWorkflow       chan *EventWorkflow
	ChanWorkflowErr    chan *EventWorkflowError
	ChanTask           chan *EventTask
	ChanTaskErr        chan *EventTaskError
	ChanSystem         chan *EventSystem
	ChanSystemErr      chan *EventSystemError
}

func NewTaskResult(t *Task) *TaskResult {
	return &TaskResult{
		TransactionID: t.TransactionID,
		TaskID:        t.TaskID,
		Status:        t.Status,
		Output:        t.Output,
		Logs:          []interface{}{},
		DoNotRetry:    false,
	}
}

func NewKafkaClient(kafkaServers string, namespace string, kafkaVersion string) (*KafkaClient, goerror.Error) {
	config := sarama.NewConfig()

	kfv, err := sarama.ParseKafkaVersion(kafkaVersion) // kafkaVersion is the version of kafka server like 0.11.0.2
	if err != nil {
		return nil, ErrUnableParseKafkaVersion
	}
	config.Version = kfv

	config.Net.MaxOpenRequests = 1

	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.Retention = 10000000 * time.Minute
	config.Consumer.Offsets.AutoCommit.Enable = true // https://github.com/Shopify/sarama/issues/1221
	config.Consumer.Fetch.Max = 100
	config.Consumer.MaxWaitTime = 100 * time.Millisecond

	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Idempotent = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 100
	config.Producer.Flush.MaxMessages = 100
	config.Producer.Flush.Frequency = time.Millisecond
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	ks := strings.Split(kafkaServers, ",")
	pd, err := sarama.NewAsyncProducer(ks, config)
	if err != nil {
		return nil, ErrUnableToCreateProducer.WithCause(err)
	}

	go func() {
		for {
			select {
			case result := <-pd.Successes():
				log.Printf("> sent to partition  %d at offset %d\n", result.Partition, result.Offset)
			case err := <-pd.Errors():
				log.Println("Failed to produce message", err)
			}
		}
	}()

	return &KafkaClient{namespace, ks, config, pd}, nil
}

func (w *KafkaClient) NewWorker(taskName string, tcb func(t *Task) *TaskResult, ccb func(t *Task) *TaskResult) goerror.Error {
	c, err := sarama.NewConsumerGroup(w.kafkaServers, fmt.Sprintf(`melonade-%s-client-%s`, w.namespace, taskName), w.config)
	if err != nil {
		return ErrUnableToCreateConsumer.WithCause(err)
	}

	wh := workerHandler{w, tcb, ccb}
	ctx := context.Background()
	go func() {
		for {
			err := c.Consume(ctx, []string{fmt.Sprintf(`melonade.%s.task.%s`, w.namespace, taskName)}, &wh)
			if err != nil {
				fmt.Printf("consume error: %v\n", err)
				time.Sleep(500 * time.Millisecond) // To prevent cpu 100% while kafka brokers are dead
			}
		}
	}()
	return nil
}

func (w *KafkaClient) NewEventWatcher(serviceName string) (*Watcher, goerror.Error) {
	c, err := sarama.NewConsumerGroup(w.kafkaServers,
		fmt.Sprintf(`melonade-%s-event-watcher-%s`, w.namespace, serviceName), w.config)
	if err != nil {
		return nil, ErrUnableToCreateConsumer.WithCause(err)
	}

	wh := eventWatcherHandler{
		&Watcher{
			make(chan *EventTransaction),
			make(chan *EventTransactionError),
			make(chan *EventWorkflow),
			make(chan *EventWorkflowError),
			make(chan *EventTask),
			make(chan *EventTaskError),
			make(chan *EventSystem),
			make(chan *EventSystemError),
		},
		w,
	}
	ctx := context.Background()
	go func() {
		for {
			err := c.Consume(ctx, []string{fmt.Sprintf(`melonade.%s.store`, w.namespace)}, &wh)
			if err != nil {
				fmt.Printf("consume error: %v\n", err)
				time.Sleep(500 * time.Millisecond) // To prevent cpu 100% while kafka brokers are dead
			}
		}
	}()
	return wh.Watcher, nil
}

// Async update task
func (w *KafkaClient) UpdateTask(tr *TaskResult) {
	trs, _ := json.Marshal(tr)

	w.producer.Input() <- &sarama.ProducerMessage{
		Topic: fmt.Sprintf(`melonade.%s.event`, w.namespace),
		Key:   sarama.StringEncoder(tr.TransactionID),
		Value: sarama.ByteEncoder(trs),
	}
}

// Async start transaction
func (w *KafkaClient) StartTransaction(tID string, wName string, wRev string, input interface{}, tags []string) {
	c := commandStartTransaction{
		TransactionID: tID,
		Type:          CommandTypeStartTransaction,
		WorkflowRef: &WorkflowRef{
			Name: wName,
			Rev:  wRev,
		},
		Input: input,
		Tags:  tags,
	}

	trs, _ := json.Marshal(c)

	w.producer.Input() <- &sarama.ProducerMessage{
		Topic: fmt.Sprintf(`melonade.%s.command`, w.namespace),
		Key:   sarama.StringEncoder(c.TransactionID),
		Value: sarama.ByteEncoder(trs),
	}
}

// workerHandler
type workerHandler struct {
	w                  *KafkaClient
	taskCallback       func(t *Task) *TaskResult
	compensateCallback func(t *Task) *TaskResult
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (wh *workerHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (wh *workerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (wh *workerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var wg sync.WaitGroup
	for m := range claim.Messages() {
		wg.Add(1)

		go func(m *sarama.ConsumerMessage) {
			defer wg.Done()
			currentMillis := time.Now().UnixNano() / int64(time.Millisecond)
			t := Task{}
			if err := json.Unmarshal(m.Value, &t); err != nil {
				fmt.Println("Bad payload", err, m.Topic, m.Partition, m.Offset)
				return
			}

			elapsedTime := int(currentMillis - t.StartTime)
			if t.AckTimeout > 0 && t.AckTimeout < elapsedTime || t.Timeout > 0 && t.Timeout < elapsedTime {
				fmt.Printf(`Already timeout: %s`, t.TaskID)
				return
			}

			tr := NewTaskResult(&t)
			tr.Status = TaskStatusInProgress
			wh.w.UpdateTask(tr)

			defer func() {
				if err := recover(); err != nil {
					tr.Status = TaskStatusFailed
					tr.Output = fmt.Sprintf("callback error: %v", err)
					wh.w.UpdateTask(tr)
				}
			}()
			switch t.Type {
			case TaskTypeTask:
				tr = wh.taskCallback(&t)
				wh.w.UpdateTask(tr)
			case TaskTypeCompensate:
				tr = wh.compensateCallback(&t)
				wh.w.UpdateTask(tr)
			default:
				tr.Status = TaskStatusFailed
				tr.Output = "Unknown task type"
				wh.w.UpdateTask(tr)
			}
			session.MarkMessage(m, "")
		}(m)
	}
	wg.Wait() // Wait for every task ran before pull for new batch
	return nil
}

// eventWatcherHandler
type eventWatcherHandler struct {
	*Watcher
	w *KafkaClient
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (eh *eventWatcherHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (eh *eventWatcherHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (eh *eventWatcherHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for m := range claim.Messages() {
		go func(m *sarama.ConsumerMessage) {
			session.MarkMessage(m, "")

			var be BaseEvent
			err := json.Unmarshal(m.Value, &be)
			if err != nil {
				log.Println(err)
				return
			}

			if be.IsError == true {
				switch be.Type {
				case EventTypeTransaction:
					var e EventTransactionError
					err := json.Unmarshal(m.Value, &e)
					if err != nil {
						log.Println(err)
						return
					}
					eh.ChanTransactionErr <- &e
				case EventTypeWorkflow:
					var e EventWorkflowError
					err := json.Unmarshal(m.Value, &e)
					if err != nil {
						log.Println(err)
						return
					}
					eh.ChanWorkflowErr <- &e
				case EventTypeTask:
					var e EventTaskError
					err := json.Unmarshal(m.Value, &e)
					if err != nil {
						log.Println(err)
						return
					}
					eh.ChanTaskErr <- &e
				case EventTypeSystem:
					var e EventSystemError
					err := json.Unmarshal(m.Value, &e)
					if err != nil {
						log.Println(err)
						return
					}
					eh.ChanSystemErr <- &e
				default:
					log.Printf(`Unknow event type: %v`, be)
				}
			} else {
				switch be.Type {
				case EventTypeTransaction:
					var e EventTransaction
					err := json.Unmarshal(m.Value, &e)
					if err != nil {
						log.Println(string(m.Value))
						log.Println(err)
						return
					}
					eh.ChanTransaction <- &e
				case EventTypeWorkflow:
					var e EventWorkflow
					err := json.Unmarshal(m.Value, &e)
					if err != nil {
						log.Println(err)
						return
					}
					eh.ChanWorkflow <- &e
				case EventTypeTask:
					var e EventTask
					err := json.Unmarshal(m.Value, &e)
					if err != nil {
						log.Println(err)
						return
					}
					eh.ChanTask <- &e
				case EventTypeSystem:
					var e EventSystem
					err := json.Unmarshal(m.Value, &e)
					if err != nil {
						log.Println(err)
						return
					}
					eh.ChanSystem <- &e
				default:
					log.Printf(`Unknow event type: %v`, be)
				}
			}
		}(m)
	}
	return nil
}
