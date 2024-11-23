Simplified orhestrator model

```mermaid
flowchart
	task(Task1)
    User --> Manager
    Manager -.-> Worker1
    Manager --> Worker2
    Manager -.-> Worker3
```

The user submits a task to the Manage node, which then selects Worker2 to run the task.
The dotted line to Worker1 and Worker3 represents that these nodes were considered but ultimately not selected to run the task. 