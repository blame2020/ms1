package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	goVersions := []string{"1.21.6"}
	arches := []string{"amd64"}

	ctx := context.Background()

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}

	defer client.Close()

	for _, ver := range goVersions {
		for _, goarch := range arches {
			err := build(ctx, client.Pipeline("build"), ver, "linux", goarch)
			if err != nil {
				panic(err)
			}
		}
	}
}

func build(ctx context.Context, client *dagger.Client, goVersion string, _ string, arch string) error {
	builder := client.Container().From(fmt.Sprintf("golang:%s", goVersion)).
		WithMountedDirectory("/src", client.Host().Directory(".")).
		WithWorkdir("/src").
		WithExec([]string{"apt-get", "update", "-y"}).
		WithExec([]string{"apt-get", "install", "-y", "make", "podman"}).
		WithEnvVariable("GOARCH", arch).
		WithExec([]string{"make", "clean"}).
		WithExec([]string{"make"}).
		WithExec([]string{"make", "docker"}).
		WithExec([]string{"make", "test"})

	if _, err := builder.Directory("build").Export(ctx, "build"); err != nil {
		return err
	}

	// if _, err := builder.Directory(".").
	// 	DockerBuild(dagger.DirectoryDockerBuildOpts{
	// 		Dockerfile: "./Dockerfile",
	// 	}).
	// 	AsTarball().Export(ctx, fmt.Sprintf("build/%s/%s/test.tar", arch)); err != nil {
	// 	return err
	// }

	return nil
}

// func save() {
// 	Publish(ctx, "test/test:latest", dagger.ContainerPublishOpts{}); err != nil {
// 		return err
// 	}

// 	return nil
// }
