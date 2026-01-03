package device

import (
	"context"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	VolumePath   = "/Volumes/NICENANO"
	PollInterval = 500 * time.Millisecond
)

// Status represents device connection status
type Status int

const (
	Disconnected Status = iota
	Connected
	Waiting
)

func (s Status) String() string {
	switch s {
	case Connected:
		return "Connected"
	case Disconnected:
		return "Disconnected"
	case Waiting:
		return "Waiting..."
	default:
		return "Unknown"
	}
}

// StatusMsg is sent when device status changes
type StatusMsg struct {
	Status    Status
	Cancelled bool
}

// Watcher polls for device connection
type Watcher struct {
	lastStatus Status
	cancelFunc context.CancelFunc
}

// NewWatcher creates a new device watcher
func NewWatcher() *Watcher {
	return &Watcher{lastStatus: Disconnected}
}

// Cancel stops any active WaitForDevice operation
func (w *Watcher) Cancel() {
	if w.cancelFunc != nil {
		w.cancelFunc()
		w.cancelFunc = nil
	}
}

// Check returns current device status
func (w *Watcher) Check() Status {
	if _, err := os.Stat(VolumePath); err == nil {
		return Connected
	}
	return Disconnected
}

// Poll returns a command that polls device status
func (w *Watcher) Poll() tea.Cmd {
	return tea.Tick(PollInterval, func(t time.Time) tea.Msg {
		status := w.Check()
		if status != w.lastStatus {
			w.lastStatus = status
			return StatusMsg{Status: status}
		}
		return nil
	})
}

// WaitForDevice returns a command that waits until device is connected
func (w *Watcher) WaitForDevice() tea.Cmd {
	// Cancel any previous wait
	w.Cancel()

	ctx, cancel := context.WithCancel(context.Background())
	w.cancelFunc = cancel

	return func() tea.Msg {
		for {
			select {
			case <-ctx.Done():
				return StatusMsg{Status: Disconnected, Cancelled: true}
			default:
				if _, err := os.Stat(VolumePath); err == nil {
					return StatusMsg{Status: Connected}
				}
				time.Sleep(PollInterval)
			}
		}
	}
}

// WaitForEject returns a command that waits until device is ejected
func WaitForEject() tea.Cmd {
	return func() tea.Msg {
		for {
			if _, err := os.Stat(VolumePath); err != nil {
				return StatusMsg{Status: Disconnected}
			}
			time.Sleep(PollInterval)
		}
	}
}
