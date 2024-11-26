package task

import "github.com/docker/go-connections/nat"

// Config - configuration for orchestration tasks

type Config struct {
	Name          string
	AttachStdin   bool
	AttachStdout  bool
	AttachStdErr  bool
	ExposedPorts  nat.PortSet
	Cmd           []string
	Image         string
	CPU           float64
	Memory        int64
	Disk          int64
	Env           []string
	RestartPolicy RestartPolicyMode
}

func NewConfig(t Task) Config {
	return Config{
		Name:          t.Name,
		Image:         t.Image,
		Memory:        t.Memory,
		Disk:          t.Disk,
		ExposedPorts:  t.ExposedPorts,
		RestartPolicy: t.RestartPolicy,
	}
}
