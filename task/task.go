package task

import (
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

type RestartPolicyMode string

const (
	RestartPolicyDisabled      RestartPolicyMode = "no"
	RestartPolicyAlways        RestartPolicyMode = "always"
	RestartPolicyOnFailure     RestartPolicyMode = "on-failure"
	RestartPolicyUnlessStopped RestartPolicyMode = "unless-stopped"
)

// Task - orchestrator task description
type Task struct {
	ID            uuid.UUID
	Name          string
	State         State
	Image         string
	Memory        int
	Disk          int
	ExposedPorts  nat.PortSet
	RestartPolicy string
	StartTime     time.Time
	FinishTime    time.Time
}

// TaskEvent - event that used by task system for state transition
type TaskEvent struct {
	ID        uuid.UUID
	State     State
	Timestamp time.Time
	Task      Task
}

// Config - configuration for orchestration tasks

type Config struct {
	Name          string
	AttachStdin   bool
	AttachStdout  bool
	AttachStdErr  bool
	ExposedPorts  nat.PortSet
	Cmd           []string
	Image         string
	Cpu           float64
	Memory        int64
	Disk          int64
	Env           []string
	RestartPolicy RestartPolicyMode
}
