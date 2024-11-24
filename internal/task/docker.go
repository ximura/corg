package task

import (
	"context"
	"io"
	"log"
	"math"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// DockerResult - result for start/stop container operations
type DockerResult struct {
	Error      error
	Action     string
	ContanerID string
	Result     string
}

// Docker - encapsulate interaction with docker SDK
type Docker struct {
	client *client.Client
}

func NewDocker() (*Docker, error) {
	dc, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &Docker{client: dc}, nil
}

// Run - run docker container
func (d *Docker) Run(ctx context.Context, config Config) (string, error) {
	reader, err := d.client.ImagePull(ctx, config.Image, image.PullOptions{})
	if err != nil {
		log.Printf("Error pulling image %s: %v\n", config.Image, err)
		return "", err
	}
	io.Copy(os.Stdout, reader)

	rp := container.RestartPolicy{
		Name: container.RestartPolicyMode(config.RestartPolicy),
	}

	r := container.Resources{
		Memory:   config.Memory,
		NanoCPUs: int64(config.CPU * math.Pow(10, 9)),
	}

	cc := container.Config{
		Image:        config.Image,
		Tty:          false,
		Env:          config.Env,
		ExposedPorts: config.ExposedPorts,
	}

	hc := container.HostConfig{
		RestartPolicy:   rp,
		Resources:       r,
		PublishAllPorts: true,
	}

	resp, err := d.client.ContainerCreate(ctx, &cc, &hc, nil, nil, config.Name)
	if err != nil {
		log.Printf("Error creating container using image %s: %v\n", config.Image, err)
		return "", err
	}
	if err := d.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Printf("Error starting container %s: %v\n", resp.ID, err)
		return "", err
	}

	out, err := d.client.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		log.Printf("Error getting logs for container %s: %v\n", resp.ID, err)
		return "", err
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return resp.ID, nil
}

func (d *Docker) Stop(ctx context.Context, id string) error {
	log.Printf("Attempting to stop container %v\n", id)
	if err := d.client.ContainerStop(ctx, id, container.StopOptions{}); err != nil {
		log.Printf("Error stopping container %s: %v\n", id, err)
		return err
	}

	if err := d.client.ContainerRemove(ctx, id, container.RemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	}); err != nil {
		log.Printf("Error removing container %s: %v\n", id, err)
		return err
	}

	return nil
}

func (d *Docker) Close() {
	d.client.Close()
}
