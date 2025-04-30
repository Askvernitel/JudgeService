package main

//TODO: REMOVE EVERY CONTAINER AND CHECK FOR io.copy errors
import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	_ "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"io"
	"time"
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
	image         string
	Stdin         io.Reader
	Stdout        io.Writer
	Stderr        io.Writer
}

func NewCmdLimiter(binPath string, memoryLimitMb int64, timeLimitSec int64) *CmdLimiter {
	return &CmdLimiter{BinPath: binPath, MemoryLimitMb: memoryLimitMb, TimeLimitSec: timeLimitSec, image: "debian:latest"}
}
func (c *CmdLimiter) pullImage(ctx context.Context, cli *client.Client) error {
	//TODO: Check If Image Exists
	_, err := cli.ImagePull(ctx, c.image, image.PullOptions{})
	if err != nil {
		return err
	}
	return err

}
func (c *CmdLimiter) createContainer(ctx context.Context, cli *client.Client) (container.CreateResponse, error) {
	resp, err := cli.ContainerCreate(ctx, &container.Config{

		Image:        c.image,
		Cmd:          []string{c.BinPath},
		Tty:          false, // stuff echoes when this is turned on OpenStdin:    true, AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		//map binary folder to docker
		Binds:     []string{"/home/danieludzlieresi/Desktop/backend-project/JudgeService/uploaded-files-tmp/:/uploaded-files-tmp"},
		Resources: container.Resources{Memory: c.MemoryLimitMb * 1024 * 1024, NanoCPUs: int64(time.Second * time.Duration(c.TimeLimitSec))},
	}, nil, nil, "")
	return resp, err

}
func (c *CmdLimiter) initContainer(ctx context.Context, cli *client.Client, resp container.CreateResponse) (types.HijackedResponse, error) {
	hijackedResp, err := cli.ContainerAttach(ctx, resp.ID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return types.HijackedResponse{}, err
	}
	return hijackedResp, err
}
func (c *CmdLimiter) Run() (*CmdResult, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return nil, err
	}
	err = c.pullImage(ctx, cli)
	if err != nil {
		return nil, err
	}
	resp, err := c.createContainer(ctx, cli)
	if err != nil {
		return nil, err
	}
	hijackedResp, err := c.initContainer(ctx, cli, resp)
	if err != nil {
		return nil, err
	}
	containerId := resp.ID
	timeOutCtx, cancel := context.WithTimeout(ctx, time.Duration(c.TimeLimitSec)*time.Second)
	var timeWroteToStdin, timeWroteToStdout time.Time
	defer cancel()
	go func() {
		_, err = io.Copy(c.Stdout, hijackedResp.Reader)
		timeWroteToStdout = time.Now()
	}()
	go func() {
		_, err = io.Copy(hijackedResp.Conn, c.Stdin)
		timeWroteToStdin = time.Now()
	}()
	defer func() {
		if err := cli.ContainerRemove(ctx, containerId, container.RemoveOptions{Force: true}); err != nil {
			fmt.Println("Container Not Removed")
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
	stdinTime := time.Since(timeWroteToStdin)
	stdoutTime := time.Since(timeWroteToStdout)
	return &CmdResult{Result: CMD_RESULT_RUN_SUCCESSFUL, TimeTakenSec: stdinTime - stdoutTime}, err
}
