package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/ximura/corg/internal/manager"
	"github.com/ximura/corg/internal/node"
	"github.com/ximura/corg/internal/task"
	"github.com/ximura/corg/internal/worker"
)

func main() {
	ctx := context.Background()
	d, err := task.NewDocker()
	if err != nil {
		log.Printf("Failed to create docker client %v", err)
		return
	}
	defer d.Close()

	fmt.Println("create a test container")
	contanerID, err := createContainer(ctx, d)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	time.Sleep(time.Second * 5)
	fmt.Printf("stopping container %s\n", contanerID)
	if err := stopContainer(ctx, d, contanerID); err != nil {
		fmt.Printf("%v", err)
		return
	}
}

func createContainer(ctx context.Context, d *task.Docker) (string, error) {
	c := task.Config{
		Name:  "test-container-1",
		Image: "postgres:13",
		Env: []string{
			"POSTGRES_USER=cube",
			"POSTGRES_PASSWORD=secret",
		},
	}

	containerID, err := d.Run(ctx, c)
	if err != nil {
		return "", err
	}

	fmt.Printf("Container %s is running with config %v\n", containerID, c)
	return containerID, nil
}

func stopContainer(ctx context.Context, d *task.Docker, id string) error {
	err := d.Stop(ctx, id)
	if err != nil {
		return err
	}

	fmt.Printf("Container %s has been stopped and removed\n", id)
	return nil
}

func printTest() {
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
		DB:    make(map[uuid.UUID]*task.Task),
	}
	fmt.Printf("worker %+v\n", w)
	w.CollectStats()
	w.RunTask()
	w.StartTask()
	w.StopTask()

	m := manager.Manager{
		Pending: *queue.New(),
		TaskDB:  make(map[string][]*task.Task),
		EventDB: make(map[string][]*task.TaskEvent),
		Workers: []string{w.Name},
	}

	fmt.Printf("manager %+v\n", m)
	m.SelectWorker()
	m.UodateTasks()
	m.SendWork()

	n := node.Node{
		Name:   "node-1",
		IP:     "192.168.1.1",
		Cores:  4,
		Memory: 1024,
		Disk:   25,
		Role:   node.Worker,
	}

	fmt.Printf("node %+v\n", n)
}
