package task

import (
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type State int

const (
	// Pending - initial state, every task started as Pending
	Pending State = iota
	// Scheduled - task moves to this state, once manager has scheduled it onto worker
	Scheduled
	// Running - task moves to this state when worker successfully starts the task
	Running
	// Completed - task moves to this state when it successfully completes its work
	Completed
	// Failed - task moves to this state when it fails completes its work
	Failed
)

type RestartPolicyMode string

const (
	RestartPolicyDisabled      RestartPolicyMode = "no"
	RestartPolicyAlways        RestartPolicyMode = "always"
	RestartPolicyOnFailure     RestartPolicyMode = "on-failure"
	RestartPolicyUnlessStopped RestartPolicyMode = "unless-stopped"
)

// TaskEvent - event that used by task system for state transition
type TaskEvent struct {
	ID        uuid.UUID
	State     State
	Timestamp time.Time
	Task      Task
}

// Task - orchestrator task description
type Task struct {
	ID            uuid.UUID
	ContainerID   string
	Name          string
	State         State
	Image         string
	Memory        int64
	Disk          int64
	ExposedPorts  nat.PortSet
	RestartPolicy RestartPolicyMode
	StartTime     time.Time
	FinishTime    time.Time
}

func (t *Task) Start(id string) {
	t.ContainerID = id
	t.State = Running
}

func (t *Task) Fail() {
	t.State = Failed
}

func (t *Task) Finish() {
	t.FinishTime = time.Now().UTC()
	t.State = Completed
}
