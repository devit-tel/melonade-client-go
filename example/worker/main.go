package worker

import (
	"flag"
	"fmt"
	melonadeClientGo "github.com/devit-tel/melonade-client-go"
	"log"
	"time"
)

func main() {
	ks := flag.String("ks", "localhost:29092", "kafka's brokers server")
	kv := flag.String("kv", "2.2.1", "kafka's version")
	ns := flag.String("ns", "docker-compose", "melonade's namespace")
	flag.Parse()

	c, err := melonadeClientGo.NewKafkaClient(*ks, *ns, *kv)
	if err != nil {
		log.Fatalln(err)
	}

	tns := []string{"t1", "t2", "t3"}
	ch := make(chan error)

	for _, tn := range tns {
		err := c.NewWorker(tn, func(t *melonadeClientGo.Task) *melonadeClientGo.TaskResult {
			fmt.Printf("Processing %s: %s: %s", t.TaskName, t.TaskID, t.TransactionID)
			time.Sleep(5 * time.Second)
			fmt.Printf("Processed %s: %s: %s", t.TaskName, t.TaskID, t.TransactionID)
			tr := t.ToTaskResult()
			tr.Status = melonadeClientGo.TaskStatusCompleted

			return tr
		}, func(t *melonadeClientGo.Task) *melonadeClientGo.TaskResult {
			fmt.Printf("Compensating %s: %s: %s", t.TaskName, t.TaskID, t.TransactionID)
			time.Sleep(5 * time.Second)
			fmt.Printf("Compensated %s: %s: %s", t.TaskName, t.TaskID, t.TransactionID)
			tr := t.ToTaskResult()
			tr.Status = melonadeClientGo.TaskStatusCompleted
			return tr
		})

		if err != nil {
			ch <- err
			log.Fatalln(err)
		}
	}

	// wait for any error
	fmt.Println(<-ch)
}
