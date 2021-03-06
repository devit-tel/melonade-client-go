// +build integration

package melonade_client_go

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

func TestWorkerClient(t *testing.T) {
	t.Run("Create worker client", func(t *testing.T) {
		const batchNumber = 120
		transactionPrefix := time.Now().Format("2006-01-02T-15-04-05")
		cli, err := NewKafkaClient("localhost:29092", "docker-compose", "2.3.1")
		if err != nil {
			t.Errorf(`should create client without error: %v`, err)
		}

		taskNames := []string{"t1", "t2", "t3"}
		var wg sync.WaitGroup
		for _, tn := range taskNames {
			err := cli.NewWorker(
				tn,
				func(t *Task) *TaskResult {
					tr := NewTaskResult(t)
					tr.Status = TaskStatusCompleted
					tr.Output = fmt.Sprintf(`Hello from %v`, t.TaskName)
					log.Printf(`Procressing %s (%s)`, t.TaskName, t.TaskID)
					return tr
				},
				func(t *Task) *TaskResult {
					tr := NewTaskResult(t)
					tr.Status = TaskStatusCompleted
					tr.Output = fmt.Sprintf(`Hello from %v`, t.TaskName)
					log.Printf(`Compensating %s (%s)`, t.TaskName, t.TaskID)
					return tr
				},
			)
			if err != nil {
				t.Errorf(`should create worker without error: %v`, err)
			}
		}

		watcher, err := cli.NewEventWatcher("local-test")
		totalTransactionEvent := batchNumber * 2
		totalWorkflowEvent := batchNumber * 2
		wg.Add(totalTransactionEvent + totalWorkflowEvent)

		go func() {
			for {
				e := <-watcher
				switch e.(type) {
				case *EventTransaction:
					td := e.(*EventTransaction)
					log.Printf("Transaction %s was %s\n", td.TransactionID, td.Details.Status)
					wg.Done()
				case *EventWorkflow:
					wd := e.(*EventWorkflow)
					log.Printf("Workflow %s was %s\n", wd.Details.WorkflowID, wd.Details.Status)
					wg.Done()
				}
			}
		}()

		time.Sleep(10 * time.Second) // Wait for kafka consumer group rebalanced
		for i := 0; i < batchNumber; i++ {
			transactionID := fmt.Sprintf(`test-%s-%d`, transactionPrefix, i)
			go cli.StartTransaction(transactionID, "simple", "1", nil, []string{"test"})
			log.Printf("Started transaction %s\n", transactionID)
		}

		wg.Wait()                   // Wait for all worker to finish (if everything work correctly)
		time.Sleep(3 * time.Second) // Wait for producer sent message
	})
}
