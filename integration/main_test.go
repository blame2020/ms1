//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/compose"
)

var stack compose.ComposeStack

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	stack, err = compose.NewDockerCompose("docker-compose.yml")

	if err != nil {
		panic(err)
	}

	if err := stack.Up(ctx, compose.Wait(true), compose.RemoveOrphans(true)); err != nil {
		panic(err)
	}

	code := m.Run()

	if err := stack.Down(ctx, compose.RemoveOrphans(true)); err != nil {
		panic(err)
	}

	os.Exit(code)
}

// ServiceIP gets a IP address of docker compose service.
func ServiceIP(t *testing.T, name string) string {
	t.Helper()

	c, err := stack.ServiceContainer(context.Background(), name)

	require.NoError(t, err)

	ip, err := c.ContainerIP(context.Background())

	require.NoError(t, err)

	return ip
}
