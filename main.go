package main

import (
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/ximura/corg/manager"
	"github.com/ximura/corg/node"
	"github.com/ximura/corg/task"
	"github.com/ximura/corg/worker"
)

func main() {
	t := task.Task{
		ID:     uuid.New(),
		Name:   "task-1",
		State:  task.Pending,
		Image:  "Image-1",
		Memory: 1024,
		Disk:   1,
	}

	te := task.TaskEvent{
		ID:    uuid.New(),
		State: task.Pending,
		Task:  t,
	}

	fmt.Printf("task %+v\n", t)
	fmt.Printf("task event %+v\n", te)

	w := worker.Worker{
		Name:  "worker-1",
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}

	fmt.Printf("worker %+v\n", w)
	w.CollectStats()
	w.RunTask()
	w.StartTask()
	w.StopTask()

	m := manager.Manager{
		Pending: *queue.New(),
		TaskDb:  make(map[string][]*task.Task),
		EventDb: make(map[string][]*task.TaskEvent),
		Workers: []string{w.Name},
	}

	fmt.Printf("manager %+v\n", m)
	m.SelectWorker()
	m.UodateTasks()
	m.SendWork()

	n := node.Node{
		Name:   "node-1",
		Ip:     "192.168.1.1",
		Cores:  4,
		Memory: 1024,
		Disk:   25,
		Role:   node.Worker,
	}

	fmt.Printf("node %+v\n", n)

	fmt.Println("create a test container")
	dockerTask, result := createContainer()
	if result.Error != nil {
		fmt.Printf("%v", result.Error)
		return
	}

	time.Sleep(time.Second * 5)
	fmt.Printf("stopping container %s\n", result.ContanerId)
	r := stopContainer(dockerTask, result.ContanerId)
	if r.Error != nil {
		fmt.Printf("%v", result.Error)
		return
	}
}

func createContainer() (*task.Docker, *task.DockerResult) {
	c := task.Config{
		Name:  "test-container-1",
		Image: "postgres:13",
		Env: []string{
			"POSTGRES_USER=cube",
			"POSTGRES_PASSWORD=secret",
		},
	}

	dc, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Println("Failed to create docker client")
		return nil, &task.DockerResult{
			Error: err,
		}
	}

	d := task.Docker{
		Client: dc,
		Config: c,
	}

	result := d.Run()
	if result.Error != nil {
		fmt.Printf("Run: %v\n", result.Error)
		return nil, nil
	}

	fmt.Printf("Container %s is running with config %v\n", result.ContanerId, c)
	return &d, &result
}

func stopContainer(d *task.Docker, id string) *task.DockerResult {
	result := d.Stop(id)
	if result.Error != nil {
		fmt.Printf("%v\n", result.Error)
		return nil
	}

	fmt.Printf("Container %s has been stopped and removed\n", result.ContanerId)
	return &result
}
