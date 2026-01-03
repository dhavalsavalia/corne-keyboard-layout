package firmware

import (
	"io"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhavalsavalia/corne-flash/internal/device"
)

// FlashStep represents what's being flashed
type FlashStep int

const (
	FlashLeft FlashStep = iota
	FlashRight
	FlashResetLeft
	FlashResetRight
)

func (s FlashStep) String() string {
	switch s {
	case FlashLeft:
		return "LEFT half"
	case FlashRight:
		return "RIGHT half"
	case FlashResetLeft:
		return "LEFT half (reset)"
	case FlashResetRight:
		return "RIGHT half (reset)"
	default:
		return "Unknown"
	}
}

// FlashProgressMsg reports flash progress
type FlashProgressMsg struct {
	Step    FlashStep
	Percent int
	Status  string
}

// FlashCompleteMsg indicates flash finished
type FlashCompleteMsg struct {
	Step    FlashStep
	Success bool
	Error   error
}

// Flasher handles firmware flashing
type Flasher struct {
	volumePath string
}

// NewFlasher creates a new flasher
func NewFlasher() *Flasher {
	return &Flasher{
		volumePath: device.VolumePath,
	}
}

// Flash copies a firmware file to the device
func (f *Flasher) Flash(step FlashStep, filePath string) tea.Cmd {
	return func() tea.Msg {
		// Open source file
		src, err := os.Open(filePath)
		if err != nil {
			return FlashCompleteMsg{Step: step, Success: false, Error: err}
		}
		defer src.Close()

		// Get file info for size
		info, err := src.Stat()
		if err != nil {
			return FlashCompleteMsg{Step: step, Success: false, Error: err}
		}

		// Create destination file
		dstPath := filepath.Join(f.volumePath, filepath.Base(filePath))
		dst, err := os.Create(dstPath)
		if err != nil {
			return FlashCompleteMsg{Step: step, Success: false, Error: err}
		}
		defer dst.Close()

		// Copy file
		written, err := io.Copy(dst, src)
		if err != nil {
			return FlashCompleteMsg{Step: step, Success: false, Error: err}
		}

		if written != info.Size() {
			return FlashCompleteMsg{
				Step:    step,
				Success: false,
				Error:   io.ErrShortWrite,
			}
		}

		return FlashCompleteMsg{Step: step, Success: true}
	}
}

