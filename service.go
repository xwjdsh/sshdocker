package dockerssh

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	LABEL_VALUE = "dockerssh"
)

var (
	cli       *client.Client
	ImageName = "dockerssh-sshd:latest"
)

func init() {
	var err error
	cli, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
}

// Create docker sshd shell
func Create(o *Options) error {
	ctx := context.Background()
	if err := build(o.Verbose, ctx); err != nil {
		return fmt.Errorf("Image build error: %v", err)
	}
	id, err := create(o, ctx)
	if err != nil {
		return fmt.Errorf("Container create error: %v", err)
	}
	if err := start(id, ctx); err != nil {
		return fmt.Errorf("Container start error: %v", err)
	}
	return nil
}

func pull() error {
	r, err := cli.ImagePull(context.Background(), "", types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, r)
	return nil
}

func build(verbose bool, ctxs ...context.Context) error {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	if err := tw.WriteHeader(&tar.Header{
		Name: "Dockerfile",
		Size: int64(len(_DOCKERFILE)),
	}); err != nil {
		return err
	}
	if _, err := tw.Write(_DOCKERFILE); err != nil {
		return err
	}

	buildOptions := types.ImageBuildOptions{
		Tags: []string{ImageName},
	}
	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	} else {
		ctx = context.Background()
	}
	resp, err := cli.ImageBuild(ctx, bytes.NewReader(buf.Bytes()), buildOptions)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	dst := ioutil.Discard
	if verbose {
		dst = os.Stdout
	}
	_, err = io.Copy(dst, resp.Body)
	return err
}

func create(o *Options, ctxs ...context.Context) (string, error) {
	containerConfig := &container.Config{
		Image:    ImageName,
		Hostname: o.Name,
		Labels:   map[string]string{"org.label-schema.name": LABEL_VALUE},
		ExposedPorts: nat.PortSet{
			"22/tcp": struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"22/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: o.Port,
				},
			},
		},
		RestartPolicy: container.RestartPolicy{Name: "always"},
	}
	if o.Volume != "" {
		hostConfig.Mounts = []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: o.Volume,
				Target: "/mnt",
			},
		}
	}

	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	} else {
		ctx = context.Background()
	}
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, o.Name)
	if err != nil {
		return "", err
	}
	return resp.ID, err
}

func start(id string, ctxs ...context.Context) error {
	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	} else {
		ctx = context.Background()
	}
	return cli.ContainerStart(ctx, id, types.ContainerStartOptions{})
}

// List all client services
func List() ([]map[string]string, error) {
	//cli, err := client.NewEnvClient()
	//if err != nil {
	//return nil, err
	//}
	//filters := filters.NewArgs()
	//filters.Add("label", fmt.Sprintf("%s=%s", LABEL_KEY, LABEL_VALUE))
	//containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
	//All:     true,
	//Filters: filters,
	//})
	//if err != nil {
	//return nil, err
	//}
	//cs := []map[string]string{}
	//for _, container := range containers {
	//name := container.Names[0]
	//if len(name) > 0 {
	//name = name[1:]
	//}
	//c := map[string]string{
	//"name":    name,
	//"state":   container.State,
	//"service": container.Labels["service"],
	//}
	//cs = append(cs, c)
	//}
	//return cs, err
	return nil, nil
}
