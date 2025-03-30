package main

//TODO: REMOVE EVERY CONTAINER AND CHECK FOR io.copy errors
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

const (
	CMD_RESULT_RUN_SUCCESSFUL        = 1
	CMD_RESULT_TIME_EXCEEDED_LIMIT   = 2
	CMD_RESULT_MEMORY_EXCEEDED_LIMIT = 3
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

func (c *CmdLimiter) Run() (*CmdResult, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return nil, err
	}
	//TODO: Check If Image Exists
	_image := "debian:latest"

	_, err = cli.ImagePull(ctx, _image, image.PullOptions{})
	if err != nil {
		return nil, err
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{

		Image:        _image,
		Cmd:          []string{c.BinPath},
		Tty:          false, // stuff echoes when this is turned on
		OpenStdin:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		//map binary folder to docker
		Binds:     []string{"/home/danieludzlieresi/Desktop/backend-project/JudgeService/uploaded-files-tmp/:/uploaded-files-tmp"},
		Resources: container.Resources{Memory: c.MemoryLimitMb * 1024 * 1024, NanoCPUs: int64(time.Second * time.Duration(c.TimeLimitSec))},
	}, nil, nil, "")
	if err != nil {
		return nil, err
	}
	containerId := resp.ID
	hijackedResp, err := cli.ContainerAttach(ctx, containerId, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err := cli.ContainerStart(ctx, containerId, container.StartOptions{}); err != nil {
		return nil, err
	}
	timeOutCtx, cancel := context.WithTimeout(ctx, time.Duration(c.TimeLimitSec)*time.Second)
	defer cancel()
	go func() {
		_, err = io.Copy(c.Stdout, hijackedResp.Reader)
	}()
	go func() {
		_, err = io.Copy(hijackedResp.Conn, c.Stdin)
	}()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cli.ContainerRemove(ctx, containerId, container.RemoveOptions{Force: true}); err != nil {
		}
	}()
	statusCh, errCh := cli.ContainerWait(ctx, containerId, container.WaitConditionNotRunning)
	select {
	case <-timeOutCtx.Done():
		if err := cli.ContainerKill(ctx, containerId, "SIGKILL"); err != nil {
			return nil, err
		}
		return &CmdResult{Result: CMD_RESULT_TIME_EXCEEDED_LIMIT}, nil
	case err = <-errCh:
		return nil, err
	case exitStatus := <-statusCh:
		fmt.Printf("Exit Code: %v", exitStatus.StatusCode)
	}
	return &CmdResult{Result: CMD_RESULT_RUN_SUCCESSFUL}, err
}
