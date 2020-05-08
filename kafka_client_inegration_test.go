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
		const batchNumber = 1000
		transactionPrefix := time.Now().Format("2006-01-02T-15-04-05")
		cli, err := NewKafkaClient("localhost:29092", "docker-compose", "2.3.1")
		if err != nil {
			t.Errorf(`should create client without error: %v`, err)
		}

		taskNames := []string{"t1", "t2", "t3"}
		var wg sync.WaitGroup
		wg.Add(len(taskNames) * batchNumber)
		for _, tn := range taskNames {
			err := cli.NewWorker(
				tn,
				func(t *Task) *taskResult {
					defer wg.Done()
					tr := NewTaskResult(t)
					tr.Status = TaskStatusCompleted
					tr.Output = fmt.Sprintf(`Hello from %v`, t.TaskName)
					log.Printf(`Procressing %s (%s)`, t.TaskName, t.TaskID)
					return tr
				},
				func(t *Task) *taskResult {
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

		time.Sleep(5 * time.Second) // Wait for kafka consumer group rebalanced
		for i := 0; i < batchNumber; i++ {
			transactionID := fmt.Sprintf(`test-%s-%d`, transactionPrefix, i)
			go cli.StartTransaction(transactionID, "simple", "1", nil, []string{"test"})
			log.Printf("Started transaction %s\n", transactionID)
		}
		wg.Wait()                   // Wait for all worker to finish (if everything work correctly)
		time.Sleep(1 * time.Second) // Wait for producer sent message
	})
}
