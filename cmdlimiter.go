package main

import (
	"context"
	"fmt"
	"io"
	"time"

	_ "github.com/docker/docker/api/types"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type CmdLimiter struct {
	BinPath       string
	TimeLimitSec  int64
	MemoryLimitMb int64
	Stdin         io.Reader
	Stdout        io.Writer
	Stderr        io.Writer
}

func NewCmdLimiter(binPath string, memoryLimitMb int64, timeLimitSec int64) *CmdLimiter {
	return &CmdLimiter{BinPath: binPath, MemoryLimitMb: memoryLimitMb, TimeLimitSec: timeLimitSec}
}

// TODO: fix this
func (c *CmdLimiter) Run() error {
	//_errCh := make(chan error)
	//	inout := make(chan []byte)
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return err
	}
	//TODO: Check If Image Exists
	_image := "debian:latest"

	_, err = cli.ImagePull(ctx, _image, image.PullOptions{})
	if err != nil {
		return err
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{

		Image:        _image,
		Cmd:          []string{c.BinPath},
		Tty:          false, // stuff is echoes when this is turned on
		OpenStdin:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		Binds:     []string{"/home/danieludzlieresi/Desktop/backend-project/JudgeService/uploaded-files-tmp/:/uploaded-files-tmp"},
		Resources: container.Resources{Memory: c.MemoryLimitMb * 1024 * 1024, NanoCPUs: int64(time.Second * time.Duration(c.TimeLimitSec))},
	}, nil, nil, "")
	if err != nil {
		return err
	}
	containerId := resp.ID
	hijackedResp, err := cli.ContainerAttach(ctx, containerId, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err := cli.ContainerStart(ctx, containerId, container.StartOptions{}); err != nil {
		return err
	}

	timeOutCtx, cancel := context.WithTimeout(ctx, time.Duration(c.TimeLimitSec)*time.Second)
	defer cancel()
	go io.Copy(c.Stdout, hijackedResp.Reader)
	go io.Copy(hijackedResp.Conn, c.Stdin)
	if err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, containerId, container.WaitConditionNotRunning)
	select {
	case <-timeOutCtx.Done():
		if err := cli.ContainerKill(ctx, containerId, "SIGKILL"); err != nil {
			return err
		}
		return fmt.Errorf("TLE")
	case err = <-errCh:
		return err
	case exitStatus := <-statusCh:
		fmt.Printf("Exit Code: %v", exitStatus.StatusCode)
	}
	if err := cli.ContainerRemove(ctx, containerId, container.RemoveOptions{Force: true}); err != nil {
		return err
	}
	return nil
}
