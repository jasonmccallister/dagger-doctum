// A generated module for DaggerDoctum functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/dagger-doctum/internal/dagger"
	"time"
)

const (
	DoctumURL  = "https://doctum.long-term.support/releases/latest/doctum.phar"
	ConfigFile = `<?php

return new Doctum\Doctum('/work/repository/');`
)

type DaggerDoctum struct{}

func (m *DaggerDoctum) Run(
	ctx context.Context,
	src *dagger.Directory,
) *dagger.Directory {
	return dag.Container(dagger.ContainerOpts{Platform: "linux/amd64"}).
		From("php:8.2-cli-alpine").
		WithWorkdir("/work").
		WithDirectory("/work/repository", src).
		WithFile("/work/bin/doctum.phar", m.GetFile(), dagger.ContainerWithFileOpts{
			Permissions: 0755,
		}).
		WithNewFile("/work/config/doctum.php", ConfigFile).
		WithWorkdir("/work/docs").
		WithEnvVariable("CACHE_BUST", time.Now().Format(time.RFC3339)).
		WithExec([]string{"/work/bin/doctum.phar", "update", "/work/config/doctum.php", "-v"}).
		Directory("/work/docs")
}

func (m *DaggerDoctum) GetFile() *dagger.File {
	return dag.Container().From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec([]string{"curl", "-O", DoctumURL}).File("doctum.phar")
}

func (m *DaggerDoctum) Serve(ctx context.Context, src *dagger.Directory) *dagger.Service {
	dir := m.Run(context.Background(), src)

	// build an nginx service
	return dag.Container().From("nginx:alpine").
		WithDirectory("/usr/share/nginx/html", dir.Directory("build")).
		WithDirectory("/usr/share/nginx/", dir.Directory("cache")).
		Terminal().
		WithExposedPort(8000).
		AsService()
}
