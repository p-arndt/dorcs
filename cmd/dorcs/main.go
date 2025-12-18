package main

import (
	"embed"
	"os"

	"github.com/p-arndt/dorcs/internal/commands"
)

//go:embed web/templates/*.html web/templates/partials/*.html
var templatesFS embed.FS

//go:embed web/static/**
var staticFS embed.FS

// Version is the version identifier for dorcs.
// This can be set at build time using -ldflags:
//
//	go build -ldflags "-X github.com/p-arndt/dorcs/cmd/dorcs.Version=1.0.0"
var Version = "dev"

func main() {
	// Check if init command is requested
	if len(os.Args) > 1 && os.Args[1] == "init" {
		commands.RunInit()
		return
	}

	// Check if build command is requested
	if len(os.Args) > 1 && os.Args[1] == "build" {
		commands.RunBuild(templatesFS, staticFS)
		return
	}

	// Default: run server mode
	commands.RunServer(templatesFS, staticFS, Version)
}
