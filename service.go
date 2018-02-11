package sshdocker

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
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	LABEL_KEY   = "org.label-schema.url"
	LABEL_VALUE = "https://github.com/xwjdsh/sshdocker"
)

var (
	cli       *client.Client
	ImageName = "sshdocker:latest"
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

func build(verbose bool, ctx context.Context) error {
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

func create(o *Options, ctx context.Context) (string, error) {
	containerConfig := &container.Config{
		Image:    ImageName,
		Hostname: o.Name,
		Labels:   map[string]string{LABEL_KEY: LABEL_VALUE},
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

	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, o.Name)
	if err != nil {
		return "", err
	}
	return resp.ID, err
}

func start(id string, ctx context.Context) error {
	return cli.ContainerStart(ctx, id, types.ContainerStartOptions{})
}

// List all docker sshd services
func List() ([]*Service, error) {
	services := []*Service{}
	containers, err := list(context.Background())
	if err != nil {
		return nil, fmt.Errorf("List services error: %v", err)
	}
	for _, container := range containers {
		s := &Service{
			Name:  container.Names[0][1:],
			State: container.State,
		}
		if len(container.Ports) > 0 {
			s.Connect = fmt.Sprintf("ssh -p %d root@localhost", container.Ports[0].PublicPort)
		}
		if len(container.Mounts) > 0 {
			s.Volume = fmt.Sprintf("%s -> %s", container.Mounts[0].Source, "/mnt")
		}
		services = append(services, s)
	}
	return services, nil
}

func list(ctx context.Context, names ...string) ([]types.Container, error) {
	filters := filters.NewArgs()
	filters.Add("label", fmt.Sprintf("%s=%s", LABEL_KEY, LABEL_VALUE))
	if len(names) > 0 {
		filters.Add("name", names[0])
	}
	return cli.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters,
	})
}

func Destroy(removeVolume bool, names []string) ([]string, []error) {
	removed, failed := []string{}, []error{}
	ctx := context.Background()
	for _, name := range names {
		cs, err := list(ctx, name)
		if err != nil {
			failed = append(failed, err)
			continue
		}
		var source string
		if len(cs) > 0 && len(cs[0].Mounts) > 0 {
			source = cs[0].Mounts[0].Source
		}
		if err := destroy(ctx, name); err == nil {
			removed = append(removed, name)
			if source != "" {
				os.RemoveAll(source)
			}
		} else {
			failed = append(failed, err)
		}
	}
	return removed, failed
}

func destroy(ctx context.Context, name string) error {
	return cli.ContainerRemove(ctx, name,
		types.ContainerRemoveOptions{
			Force: true,
		},
	)
}
