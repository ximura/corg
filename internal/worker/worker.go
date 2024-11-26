package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/ximura/corg/internal/task"
)

// Worker - responsible for performing work in form of Tasks
// It represent physical or virual machine that can execute Task
type Worker struct {
	Docker    *task.Docker
	Name      string
	Queue     queue.Queue
	DB        map[uuid.UUID]*task.Task
	TaskCount int
}

func (w *Worker) CollectStats() {
	fmt.Println("I will collect stats")
}
func (w *Worker) AddTask(t task.Task) {
	w.Queue.Enqueue(t)
}

func (w *Worker) RunTask(ctx context.Context) (string, error) {
	t := w.Queue.Dequeue()
	if t == nil {
		log.Println("No tasks in the queue")
		return "", nil
	}

	taskQueued := t.(task.Task)
	taskPersisted := w.DB[taskQueued.ID]
	// new task need to add to storage
	if taskPersisted == nil {
		taskPersisted = &taskQueued
		w.DB[taskQueued.ID] = &taskQueued
	}

	if task.ValidStateTransition(taskPersisted.State, taskQueued.State) {
		switch taskQueued.State {
		case task.Scheduled:
			return w.StartTask(ctx, taskQueued)
		case task.Completed:
			return "", w.StopTask(ctx, taskQueued)
		default:
			return "", fmt.Errorf("Unsupported state %d", taskQueued.State)
		}
	}

	return "", fmt.Errorf("Invalid transition from %v to %v", taskPersisted.State, taskQueued.State)
}

func (w *Worker) StartTask(ctx context.Context, t task.Task) (string, error) {
	config := task.NewConfig(t)
	t.StartTime = time.Now().UTC()
	containerID, err := w.Docker.Run(ctx, config)
	if err != nil {
		log.Printf("Error running task %v: %v\n", t.ID, err)
		t.Fail()
		w.DB[t.ID] = &t
		return "", err
	}
	t.Start(containerID)
	w.DB[t.ID] = &t

	return containerID, nil
}

func (w *Worker) StopTask(ctx context.Context, t task.Task) error {
	if err := w.Docker.Stop(ctx, t.ContainerID); err != nil {
		log.Printf("Error stopping container %v: %v\n", t.ContainerID, err)
		return err
	}

	t.Finish()
	w.DB[t.ID] = &t
	log.Printf("Stopped and removed container %v for task %v\n", t.ContainerID, t.ID)

	return nil
}
