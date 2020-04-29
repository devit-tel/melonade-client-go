melonade-client-go
---

- [x] start workflow
- [ ] worker client


Usage
---
```go
	client := New([PROCESS_MANAGER_ENDPOINT])
        resp, _ := client.StartWorkflow([TASK_NAME], [REVISION], [TRANSACTION_ID], [PAYLOAD_OBJECT])
```
