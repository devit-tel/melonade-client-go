package melonade_client_go

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/devit-tel/goerror"
)

type kafkaClient struct {
	namespace    string
	kafkaServers []string
	config       *sarama.Config
	producer     sarama.AsyncProducer
}

func NewTaskResult(t *Task) *taskResult {
	return &taskResult{
		TransactionID: t.TransactionID,
		TaskID:        t.TaskID,
		Status:        t.Status,
		Output:        t.Output,
		Logs:          []interface{}{},
		DoNotRetry:    false,
	}
}

func NewWorkerClient(kafkaServers string, namespace string, kafkaVersion string) (*kafkaClient, goerror.Error) {
	config := sarama.NewConfig()

	kfv, err := sarama.ParseKafkaVersion(kafkaVersion) // kafkaVersion is the version of kafka server like 0.11.0.2
	if err != nil {
		return nil, ErrUnableParseKafkaVersion
	}
	config.Version = kfv

	config.Net.MaxOpenRequests = 1

	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Fetch.Max = 100

	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Idempotent = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10000000
	config.Producer.Flush.MaxMessages = 100
	config.Producer.Flush.Frequency = time.Millisecond

	ks := strings.Split(kafkaServers, ",")
	pd, err := sarama.NewAsyncProducer(ks, config)
	if err != nil {
		return nil, goerror.DefineInternalServerError("UnableToCreateProducer", fmt.Sprint(err))
	}

	return &kafkaClient{namespace, ks, config, pd}, nil
}

func (w *kafkaClient) NewWorker(taskName string, tcb func(t *Task) *taskResult, ccb func(t *Task) *taskResult) goerror.Error {
	c, err := sarama.NewConsumerGroup(w.kafkaServers, fmt.Sprintf(`melonade-%s.client`, w.namespace), w.config)
	if err != nil {
		return goerror.DefineInternalServerError("UnableToCreateConsumer", fmt.Sprint(err))
	}

	wh := WorkerHandler{w, tcb, ccb}
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

// Async update task
func (w *kafkaClient) UpdateTask(tr *taskResult) {
	trs, _ := json.Marshal(tr)

	w.producer.Input() <- &sarama.ProducerMessage{
		Topic: fmt.Sprintf(`melonade.%s.event`, w.namespace),
		Key:   sarama.StringEncoder(tr.TransactionID),
		Value: sarama.ByteEncoder(trs),
	}
}

// Async start transaction
func (w *kafkaClient) StartTransaction(tID string, wName string, wRev string, input interface{}, tags []string) {
	c := commandStartTransaction{
		TransactionID: tID,
		Type:          CommandTypeStartTransaction,
		WorkflowRef: &workflowRef{
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

// WorkerHandler
type WorkerHandler struct {
	w                  *kafkaClient
	taskCallback       func(t *Task) *taskResult
	compensateCallback func(t *Task) *taskResult
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (wh *WorkerHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (wh *WorkerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (wh *WorkerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
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
