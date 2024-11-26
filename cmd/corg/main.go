package main

import (
	"context"
	"fmt"
	"sync"
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
		fmt.Printf("Failed to create docker client %v", err)
		return
	}
	defer d.Close()

	workers := []worker.Worker{
		createWorker(d), createWorker(d), createWorker(d),
	}

	var wg sync.WaitGroup
	wg.Add(3)
	for i, w := range workers {
		go func() {
			defer wg.Done()
			t := task.Task{
				ID:    uuid.New(),
				Name:  fmt.Sprintf("test-container-%d", i),
				State: task.Scheduled,
				Image: "strm/helloworld-http",
			}

			fmt.Printf("starting task %s\n", t.ID)
			w.AddTask(t)
			containerID, err := w.RunTask(ctx)
			if err != nil {
				fmt.Printf("Failed to run task %s. %v\n", t.ID, err)
				return
			}

			t.ContainerID = containerID
			fmt.Printf("task %s is running in container %s\n", t.ID, t.ContainerID)
			fmt.Println("Sleepy time")
			time.Sleep(30 * time.Second)

			fmt.Printf("stopping task %s\n", t.ID)
			t.State = task.Completed
			w.AddTask(t)
			if _, err := w.RunTask(ctx); err != nil {
				fmt.Printf("Failed to run task %s. %v\n", t.ID, err)
				return
			}
		}()
	}

	wg.Wait()
}

func createWorker(d *task.Docker) worker.Worker {
	return worker.Worker{
		Queue:  *queue.New(),
		DB:     make(map[uuid.UUID]*task.Task),
		Docker: d,
	}
}

func testTask() error {
	ctx := context.Background()
	d, err := task.NewDocker()
	if err != nil {
		return fmt.Errorf("Failed to create docker client %v", err)
	}
	defer d.Close()

	fmt.Println("create a test container")
	contanerID, err := createContainer(ctx, d)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 5)
	fmt.Printf("stopping container %s\n", contanerID)
	if err := stopContainer(ctx, d, contanerID); err != nil {
		return err
	}
	return nil
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

func printTest(ctx context.Context) {
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
	w.RunTask(ctx)
	w.StartTask(ctx, t)
	w.StopTask(ctx, t)

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
