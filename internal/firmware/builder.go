package firmware

import (
	"bufio"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// BuildStep represents a build step
type BuildStep int

const (
	StepActivateEnv BuildStep = iota
	StepBuildLeft
	StepBuildRight
	StepCopyFirmware
	StepComplete
)

func (s BuildStep) String() string {
	switch s {
	case StepActivateEnv:
		return "Activate environment"
	case StepBuildLeft:
		return "Build left half"
	case StepBuildRight:
		return "Build right half"
	case StepCopyFirmware:
		return "Copy to firmware/"
	case StepComplete:
		return "Complete"
	default:
		return "Unknown"
	}
}

// BuildProgressMsg reports build progress
type BuildProgressMsg struct {
	Step    BuildStep
	Percent int
	Output  string
}

// BuildCompleteMsg indicates build finished
type BuildCompleteMsg struct {
	Success bool
	Error   error
	Path    string
}

// Builder runs firmware builds
type Builder struct {
	scriptPath string
}

// NewBuilder creates a new builder
func NewBuilder(baseDir string) *Builder {
	return &Builder{
		scriptPath: filepath.Join(baseDir, "build.sh"),
	}
}

// Build runs the build script
func (b *Builder) Build(target string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command(b.scriptPath, target)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return BuildCompleteMsg{Success: false, Error: err}
		}

		if err := cmd.Start(); err != nil {
			return BuildCompleteMsg{Success: false, Error: err}
		}

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			// Could parse output for progress, but for now just run
			_ = scanner.Text()
		}

		if err := cmd.Wait(); err != nil {
			return BuildCompleteMsg{Success: false, Error: err}
		}

		return BuildCompleteMsg{Success: true}
	}
}

// BuildWithProgress runs build and reports progress
func (b *Builder) BuildWithProgress(target string, progressChan chan<- BuildProgressMsg) tea.Cmd {
	return func() tea.Msg {
		steps := []BuildStep{StepActivateEnv, StepBuildLeft, StepBuildRight, StepCopyFirmware}

		for i, step := range steps {
			progressChan <- BuildProgressMsg{
				Step:    step,
				Percent: (i * 100) / len(steps),
			}
		}

		cmd := exec.Command(b.scriptPath, target)
		output, err := cmd.CombinedOutput()

		if err != nil {
			return BuildCompleteMsg{Success: false, Error: err}
		}

		progressChan <- BuildProgressMsg{
			Step:    StepComplete,
			Percent: 100,
			Output:  string(output),
		}

		return BuildCompleteMsg{Success: true}
	}
}
