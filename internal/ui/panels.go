package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dhavalsavalia/corne-flash/internal/firmware"
)

// Panel identifiers
type Panel int

const (
	PanelFirmware Panel = iota
	PanelStatus
	PanelLog
)

func (p Panel) String() string {
	switch p {
	case PanelFirmware:
		return "Firmware"
	case PanelStatus:
		return "Status"
	case PanelLog:
		return "Log"
	default:
		return "Unknown"
	}
}

// LogEntry represents a log message
type LogEntry struct {
	Time    time.Time
	Message string
	Level   LogLevel
}

// LogLevel for log entries
type LogLevel int

const (
	LogInfo LogLevel = iota
	LogSuccess
	LogWarning
	LogError
)

// FirmwarePanel renders the firmware list
type FirmwarePanel struct {
	builds   []firmware.Build
	selected int
	height   int
	width    int
}

// NewFirmwarePanel creates a new firmware panel
func NewFirmwarePanel() *FirmwarePanel {
	return &FirmwarePanel{
		selected: 0,
	}
}

// SetBuilds updates the firmware builds list
func (p *FirmwarePanel) SetBuilds(builds []firmware.Build) {
	p.builds = builds
	if p.selected >= len(builds) {
		p.selected = len(builds) - 1
	}
	if p.selected < 0 {
		p.selected = 0
	}
}

// Selected returns the selected build
func (p *FirmwarePanel) Selected() *firmware.Build {
	if len(p.builds) == 0 {
		return nil
	}
	return &p.builds[p.selected]
}

// MoveUp moves selection up
func (p *FirmwarePanel) MoveUp() {
	if p.selected > 0 {
		p.selected--
	}
}

// MoveDown moves selection down
func (p *FirmwarePanel) MoveDown() {
	if p.selected < len(p.builds)-1 {
		p.selected++
	}
}

// SetSize sets the panel dimensions
func (p *FirmwarePanel) SetSize(width, height int) {
	p.width = width
	p.height = height
}

// View renders the firmware panel content
func (p *FirmwarePanel) View() string {
	if len(p.builds) == 0 {
		return DimStyle.Render("  No firmware found")
	}

	var lines []string
	for i, build := range p.builds {
		prefix := "  "
		if i == p.selected {
			prefix = "▸ "
		}

		// Format date
		dateStr := firmware.FormatDate(build.Date)

		// Status indicator
		status := ""
		if build.HasLeft && build.HasRight {
			status = SuccessStyle.Render(" ✓")
		}

		line := prefix + dateStr + status
		if i == p.selected {
			line = SelectedStyle.Render(line)
		}
		lines = append(lines, line)

		// Show files for selected build
		if i == p.selected {
			for j, f := range build.Files {
				treeChr := TreeBranch
				if j == len(build.Files)-1 {
					treeChr = TreeLast
				}
				fileLine := fmt.Sprintf("  %s %s", treeChr, f.Name)
				lines = append(lines, DimStyle.Render(fileLine))
			}
		}
	}

	return strings.Join(lines, "\n")
}

// StatusPanel renders the status/operation display
type StatusPanel struct {
	width  int
	height int
}

// NewStatusPanel creates a new status panel
func NewStatusPanel() *StatusPanel {
	return &StatusPanel{}
}

// SetSize sets the panel dimensions
func (p *StatusPanel) SetSize(width, height int) {
	p.width = width
	p.height = height
}

// ViewIdle renders idle state
func (p *StatusPanel) ViewIdle(build *firmware.Build) string {
	var lines []string

	// Center box
	boxWidth := p.width - 8
	if boxWidth < 20 {
		boxWidth = 20
	}

	lines = append(lines, "")
	lines = append(lines, centerText("SELECT FIRMWARE", boxWidth))
	lines = append(lines, "")
	lines = append(lines, centerText("Choose a build to flash", boxWidth))
	lines = append(lines, centerText("or press B to build new", boxWidth))
	lines = append(lines, "")

	if build != nil {
		lines = append(lines, "")
		lines = append(lines, DimStyle.Render("Selected: ")+firmware.FormatDate(build.Date))
	}

	return strings.Join(lines, "\n")
}

// ViewBuilding renders building state
func (p *StatusPanel) ViewBuilding(step firmware.BuildStep, percent int) string {
	var lines []string

	spinner := SpinnerFrames[(time.Now().UnixMilli()/100)%int64(len(SpinnerFrames))]

	lines = append(lines, "")
	lines = append(lines, AccentStyle.Render(spinner+" BUILDING FIRMWARE"))
	lines = append(lines, "")
	lines = append(lines, RenderProgressBar(percent, p.width-10))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Step: %s", step.String()))
	lines = append(lines, "")

	// Step checklist
	steps := []firmware.BuildStep{
		firmware.StepActivateEnv,
		firmware.StepBuildLeft,
		firmware.StepBuildRight,
		firmware.StepCopyFirmware,
	}

	for _, s := range steps {
		icon := "◻"
		style := DimStyle
		if s < step {
			icon = "◼"
			style = SuccessStyle
		} else if s == step {
			icon = "◼"
			style = AccentStyle
		}
		lines = append(lines, style.Render(icon+" "+s.String()))
	}

	return strings.Join(lines, "\n")
}

// ViewWaiting renders waiting for device state
func (p *StatusPanel) ViewWaiting(target string) string {
	var lines []string

	spinner := SpinnerFrames[(time.Now().UnixMilli()/100)%int64(len(SpinnerFrames))]

	lines = append(lines, "")
	lines = append(lines, "")
	lines = append(lines, centerText(WarningStyle.Render(spinner+" WAITING FOR DEVICE"), p.width))
	lines = append(lines, "")
	lines = append(lines, centerText("Connect "+target, p.width))
	lines = append(lines, centerText("Double-tap reset", p.width))
	lines = append(lines, "")
	lines = append(lines, "")
	lines = append(lines, DimStyle.Render("Polling /Volumes/NICENANO..."))

	return strings.Join(lines, "\n")
}

// ViewFlashing renders flashing in progress
func (p *StatusPanel) ViewFlashing(step firmware.FlashStep, percent int, filename string) string {
	var lines []string

	spinner := SpinnerFrames[(time.Now().UnixMilli()/100)%int64(len(SpinnerFrames))]

	lines = append(lines, "")
	lines = append(lines, AccentStyle.Render(spinner+" FLASHING "+step.String()))
	lines = append(lines, "")
	lines = append(lines, RenderProgressBar(percent, p.width-10))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Copying: %s", filename))
	lines = append(lines, "")

	// Flash checklist
	lines = append(lines, "")
	leftIcon := "◻"
	rightIcon := "◻"
	if step == firmware.FlashLeft {
		leftIcon = "◼"
	} else if step == firmware.FlashRight {
		leftIcon = "◼"
		rightIcon = "◼"
	}
	lines = append(lines, AccentStyle.Render(leftIcon+" Flash left half"))
	lines = append(lines, DimStyle.Render(rightIcon+" Flash right half"))

	return strings.Join(lines, "\n")
}

// ViewComplete renders completion summary
func (p *StatusPanel) ViewComplete(duration time.Duration, steps []string) string {
	var lines []string

	lines = append(lines, "")
	lines = append(lines, SuccessStyle.Render("✓ FLASH COMPLETE"))
	lines = append(lines, "")

	for _, step := range steps {
		lines = append(lines, SuccessStyle.Render("  ✓  "+step))
	}

	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  Duration: %s", duration.Round(time.Second)))
	lines = append(lines, "")
	lines = append(lines, DimStyle.Render("Test both halves to verify."))

	return strings.Join(lines, "\n")
}

// LogPanel renders the log output
type LogPanel struct {
	entries []LogEntry
	width   int
	height  int
}

// NewLogPanel creates a new log panel
func NewLogPanel() *LogPanel {
	return &LogPanel{}
}

// Add adds a log entry
func (p *LogPanel) Add(level LogLevel, msg string) {
	p.entries = append(p.entries, LogEntry{
		Time:    time.Now(),
		Message: msg,
		Level:   level,
	})
	// Keep last N entries
	maxEntries := 50
	if len(p.entries) > maxEntries {
		p.entries = p.entries[len(p.entries)-maxEntries:]
	}
}

// Clear clears all entries
func (p *LogPanel) Clear() {
	p.entries = nil
}

// SetSize sets the panel dimensions
func (p *LogPanel) SetSize(width, height int) {
	p.width = width
	p.height = height
}

// View renders the log panel content
func (p *LogPanel) View() string {
	if len(p.entries) == 0 {
		return DimStyle.Render("  No log entries")
	}

	// Show entries that fit in height
	maxVisible := p.height - 2
	if maxVisible < 1 {
		maxVisible = 10
	}

	start := 0
	if len(p.entries) > maxVisible {
		start = len(p.entries) - maxVisible
	}

	var lines []string
	for _, entry := range p.entries[start:] {
		timestamp := DimStyle.Render(entry.Time.Format("15:04:05"))

		var msgStyle lipgloss.Style
		switch entry.Level {
		case LogSuccess:
			msgStyle = SuccessStyle
		case LogWarning:
			msgStyle = WarningStyle
		case LogError:
			msgStyle = ErrorStyle
		default:
			msgStyle = lipgloss.NewStyle().Foreground(ColorFg)
		}

		// Truncate message if needed
		msg := entry.Message
		maxMsgLen := p.width - 12
		if maxMsgLen > 0 && len(msg) > maxMsgLen {
			msg = msg[:maxMsgLen-3] + "..."
		}

		lines = append(lines, timestamp+"  "+msgStyle.Render(msg))
	}

	return strings.Join(lines, "\n")
}

// Helper functions
func centerText(text string, width int) string {
	textLen := lipgloss.Width(text)
	if textLen >= width {
		return text
	}
	padding := (width - textLen) / 2
	return strings.Repeat(" ", padding) + text
}
