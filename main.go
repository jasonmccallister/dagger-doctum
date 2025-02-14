// A module used to generate documentation for PHP projects.
package main

import (
	"context"
	"dagger/dagger-doctum/internal/dagger"
	"fmt"
)

const (
	DoctumURLFmt      = "https://doctum.long-term.support/releases/%s/doctum.phar"
	DefaultConfigFile = `<?php

return new Doctum\Doctum('/work/repository/');`
)

type DaggerDoctum struct {
	// +private
	Version string
	// +private
	Image string
	// +private
	ConfigFile string
	// +private
	Source *dagger.Directory
}

func New(
	// The source directory containing the PHP files to document.
	src *dagger.Directory,
	// The version of Doctum to use.
	// +optional
	// +default="5.5.4"
	version string,
	// The image to use for building the documentation.
	// +optional
	// +default="php:8.3-cli-alpine"
	image string,
	// The path to the Doctum configuration file.
	// +optional
	configFile *dagger.File,
) (*DaggerDoctum, error) {
	var configuration string
	if configFile == nil {
		configuration = DefaultConfigFile
	}

	content, err := configFile.Contents(context.Background())
	if err != nil {
		return nil, err
	}

	configuration = content

	return &DaggerDoctum{
		Version:    version,
		Image:      image,
		ConfigFile: configuration,
		Source:     src,
	}, nil
}

// Run builds the documentation for the given source directory.
func (m *DaggerDoctum) Run(ctx context.Context) *dagger.Directory {
	doctum := dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "curl"}).
		WithExec([]string{"curl", "-O", fmt.Sprintf(DoctumURLFmt, m.Version)}).
		File("doctum.phar")

	return dag.Container().
		From(m.Image).
		WithWorkdir("/work").
		WithDirectory("/work/repository", m.Source).
		WithFile("/work/bin/doctum.phar", doctum, dagger.ContainerWithFileOpts{
			Permissions: 0755,
		}).
		WithNewFile("/work/config/doctum.php", m.ConfigFile).
		WithWorkdir("/work/docs").
		WithExec([]string{"/work/bin/doctum.phar", "update", "/work/config/doctum.php", "-v"}).
		Directory("/work/docs")
}

// Serve builds the documentation for the given source directory and serves it.
func (m *DaggerDoctum) Serve(ctx context.Context) *dagger.Service {
	return dag.Container().From("nginx:alpine").
		WithDirectory("/usr/share/nginx/html", m.Run(context.Background()).Directory("build")).
		WithExposedPort(80).
		AsService()
}
