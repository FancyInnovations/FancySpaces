package runtime

import "context"

type Runtime interface {
	ListContainers(ctx context.Context) ([]Container, error)
	CreateContainer(ctx context.Context, spec ContainerSpec) (Container, error)
	DeleteContainer(ctx context.Context, id string) error
	RestartContainer(ctx context.Context, id string) error
}
