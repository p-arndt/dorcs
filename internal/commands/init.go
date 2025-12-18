package commands

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// RunInit handles the init subcommand for initializing a new documentation site.
func RunInit() {
	initFlags := flag.NewFlagSet("init", flag.ExitOnError)
	var (
		dir        = initFlags.String("dir", "./docs", "Directory to initialize (will be created if it doesn't exist)")
		title      = initFlags.String("title", "Documentation", "Site title")
		withConfig = initFlags.Bool("config", true, "Create a basic dorcs.yaml config file")
	)
	initFlags.Parse(os.Args[2:]) // Skip "init" command

	// Resolve directory path
	absDir, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatalf("resolve dir: %v", err)
	}

	// Check if directory already exists
	if st, err := os.Stat(absDir); err == nil {
		if !st.IsDir() {
			log.Fatalf("path exists but is not a directory: %s", absDir)
		}
		// Directory exists, check if it's empty
		entries, err := os.ReadDir(absDir)
		if err != nil {
			log.Fatalf("read directory: %v", err)
		}
		if len(entries) > 0 {
			log.Printf("dorcs: warning: directory %s already exists and is not empty", absDir)
			log.Printf("dorcs: continuing anyway...")
		}
	} else {
		// Create directory
		if err := os.MkdirAll(absDir, 0755); err != nil {
			log.Fatalf("create directory: %v", err)
		}
		log.Printf("dorcs: created directory %s", absDir)
	}

	// Create index.md
	indexPath := filepath.Join(absDir, "index.md")
	if _, err := os.Stat(indexPath); err == nil {
		log.Printf("dorcs: warning: index.md already exists, skipping...")
	} else {
		// Get current date in YYYY-MM-DD format
		now := time.Now()
		dateStr := now.Format("2006-01-02")

		indexContent := fmt.Sprintf(`---
title: "%s"
description: "Welcome to %s"
tags: [docs]
date: %s
draft: false
---

# %s

Welcome to your documentation site!

## Getting Started

This is your index page. You can start editing this file to customize your documentation.

## Features

- **Easy to use** - Just write Markdown files
- **Fast** - Single binary, no dependencies
- **Beautiful** - Modern, responsive design
- **Search** - Built-in search functionality

## Next Steps

1. Edit this file (index.md) to customize your homepage
2. Add more Markdown files to create additional pages
3. Use 'dorcs --watch' to start the development server with live reload
4. Use 'dorcs build' to generate a static site for deployment

Happy documenting!
`, *title, *title, dateStr, *title)

		if err := os.WriteFile(indexPath, []byte(indexContent), 0644); err != nil {
			log.Fatalf("write index.md: %v", err)
		}
		log.Printf("dorcs: created %s", indexPath)
	}

	// Create config file if requested
	if *withConfig {
		configPath := "dorcs.yaml"
		// Check if config already exists in current directory
		if _, err := os.Stat(configPath); err == nil {
			log.Printf("dorcs: warning: %s already exists, skipping...", configPath)
		} else {
			// Check if config exists in docs directory
			configInDocs := filepath.Join(absDir, "dorcs.yaml")
			if _, err := os.Stat(configInDocs); err == nil {
				log.Printf("dorcs: warning: %s already exists, skipping...", configInDocs)
			} else {
				// Create basic config file
				configContent := fmt.Sprintf(`# dorcs Configuration File
# Place this file in your project root (current working directory) or docs directory

# Server port (default: 8080)
port: 8080

# Site metadata
site:
  # Title shown in the header/brand area (top-left)
  title: "%s"

  # Description for meta tags (optional)
  description: "Documentation site"

# Theme and styling configuration
theme:
  # Mode: "light", "dark", or "auto" (follows system preference)
  mode: auto

  # Preset theme: default, ocean, forest, sunset, midnight, lavender, rose
  preset: default

# Navigation configuration
nav:
  # Show search box in sidebar (default: true)
  show_search: true

  # Keep all folders expanded by default (default: false)
  expand_all: false

# Footer configuration
footer:
  # Custom footer text
  text: "Built with ❤️"

  # Show "Powered by dorcs" (default: true)
  show_powered_by: true
`, *title)

				if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
					log.Fatalf("write config file: %v", err)
				}
				log.Printf("dorcs: created %s", configPath)
			}
		}
	}

	log.Printf("dorcs: initialization complete!")
	log.Printf("dorcs: documentation directory: %s", absDir)
	log.Printf("dorcs: run 'dorcs --watch' to start the development server")
}
