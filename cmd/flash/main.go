package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhavalsavalia/corne-flash/internal/ui"
)

func main() {
	// Determine base directory (where this script lives)
	baseDir, err := getBaseDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Create model
	model := ui.NewModel(baseDir)

	// Run the program
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func getBaseDir() (string, error) {
	// Try to find corne-build directory
	// First, check if we're running from within it
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Check if firmware/ exists in current directory
	if _, err := os.Stat(filepath.Join(cwd, "firmware")); err == nil {
		return cwd, nil
	}

	// Check if build.sh exists in current directory
	if _, err := os.Stat(filepath.Join(cwd, "build.sh")); err == nil {
		return cwd, nil
	}

	// Try executable directory
	exe, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exe)
		if _, err := os.Stat(filepath.Join(exeDir, "firmware")); err == nil {
			return exeDir, nil
		}
	}

	// Fall back to current directory
	return cwd, nil
}
