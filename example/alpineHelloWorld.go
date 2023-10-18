package example

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"os"
	"wilikidi/docker-client/utils"
)

// equals to：docker run alpine echo hello world

func AlpineHelloWorld() {

	// 初始化一个上下文。
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// 拉取镜像。返回一个reader和一个错误。这个reader是拉取镜像时的错误。
	reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	// 这块如果分批次读取，那么读取到的信息是间断的。
	//info := make([]byte, 1024)
	//time.Sleep(time.Second)
	//fmt.Println(reader.Read(info))
	//fmt.Println(string(info))

	//fmt.Println("--------------------------->")
	//fmt.Println(reader.Read(info))
	//fmt.Println(string(info))
	// 本函数从reader读取到的字符，打印到标准输出中。
	//fmt.Println("===========================>")
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
	}, nil, nil, nil, "")

	if err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}

func RunContainerBackground() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	imageName := "bfirsh/reticulate-splines"
	reader, err := utils.CLIENT.ImagePull(imageName)
	if err != nil {
		panic(err)
	}

	// 重定向到输出
	writeNumber, err := io.Copy(os.Stdout, reader)
	if err != nil {
		panic(err)
	}
	fmt.Printf("写出的数字大小为%v\n", writeNumber)

	config := container.Config{Image: imageName}
	resp, err := utils.CLIENT.ContainerCreate(&config)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)

}

func ListManagerContainers() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container)
	}

}

func StopRunningContainers() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
			panic(err)
		}
		fmt.Println(container.ID, "is stoped")
	}
}

func PrintTheLog() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	out, err := cli.ContainerLogs(ctx, "5757d5bb221b", types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)
}

func TestPullImageWithAuthentication() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	authConfig := types.AuthConfig{
		Username: "username",
		Password: "password",
	}
	encodeJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodeJSON)
	out, err := cli.ImagePull(ctx, "docker.tvunetworks.com/tvumma/batchjob:1.0.0.48build48", types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		panic(err)
	}

	defer out.Close()
	io.Copy(os.Stdout, out)
}

func CommitContainer() {

}
