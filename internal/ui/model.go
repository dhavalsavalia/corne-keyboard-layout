package ui

import (
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dhavalsavalia/corne-flash/internal/device"
	"github.com/dhavalsavalia/corne-flash/internal/firmware"
)

// AppState represents the application state
type AppState int

const (
	StateIdle AppState = iota
	StateBuilding
	StateWaitingDevice
	StateFlashing
	StateFactoryReset
	StateComplete
)

// Model is the main bubbletea model
type Model struct {
	// Dimensions
	width  int
	height int

	// State
	state        AppState
	activePanel  Panel
	showHelp     bool
	showDialog   bool
	deviceStatus device.Status

	// Panels
	firmwarePanel *FirmwarePanel
	statusPanel   *StatusPanel
	logPanel      *LogPanel

	// Overlays
	helpOverlay   *HelpOverlay
	confirmDialog *ConfirmDialog

	// Core components
	scanner *firmware.Scanner
	watcher *device.Watcher
	builder *firmware.Builder
	flasher *firmware.Flasher

	// Operation state
	buildStep     firmware.BuildStep
	buildPercent  int
	flashStep     firmware.FlashStep
	flashPercent  int
	flashTarget   string
	startTime     time.Time
	completedSteps []string

	// Base directory
	baseDir string
}

// NewModel creates a new model
func NewModel(baseDir string) *Model {
	firmwareDir := filepath.Join(baseDir, "firmware")

	return &Model{
		baseDir:       baseDir,
		state:         StateIdle,
		activePanel:   PanelFirmware,
		deviceStatus:  device.Disconnected,
		firmwarePanel: NewFirmwarePanel(),
		statusPanel:   NewStatusPanel(),
		logPanel:      NewLogPanel(),
		helpOverlay:   NewHelpOverlay(),
		scanner:       firmware.NewScanner(firmwareDir),
		watcher:       device.NewWatcher(),
		builder:       firmware.NewBuilder(baseDir),
		flasher:       firmware.NewFlasher(),
	}
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	// Scan firmware and check device
	m.logPanel.Add(LogInfo, "App started")

	builds, err := m.scanner.Scan()
	if err != nil {
		m.logPanel.Add(LogError, "Scan failed: "+err.Error())
	} else {
		m.firmwarePanel.SetBuilds(builds)
		m.logPanel.Add(LogInfo, "Found "+formatInt(len(builds))+" builds")
	}

	// Initial device check
	m.deviceStatus = m.watcher.Check()
	if m.deviceStatus == device.Connected {
		m.logPanel.Add(LogSuccess, "Device connected")
	}

	return m.watcher.Poll()
}

// Update handles messages
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updatePanelSizes()
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case device.StatusMsg:
		if msg.Cancelled {
			// Ignore cancelled wait messages
			return m, nil
		}
		m.deviceStatus = msg.Status
		if msg.Status == device.Connected {
			m.logPanel.Add(LogSuccess, "Device connected")
			if m.state == StateWaitingDevice {
				return m.startFlash()
			}
		} else {
			m.logPanel.Add(LogInfo, "Device disconnected")
		}
		return m, m.watcher.Poll()

	case firmware.BuildCompleteMsg:
		if msg.Success {
			m.logPanel.Add(LogSuccess, "Build complete")
			m.buildPercent = 100
			// Refresh firmware list
			builds, _ := m.scanner.Scan()
			m.firmwarePanel.SetBuilds(builds)
			m.state = StateIdle
		} else {
			m.logPanel.Add(LogError, "Build failed: "+msg.Error.Error())
			m.state = StateIdle
		}
		return m, nil

	case firmware.FlashCompleteMsg:
		if msg.Success {
			m.logPanel.Add(LogSuccess, msg.Step.String()+" flashed")
			m.completedSteps = append(m.completedSteps, msg.Step.String()+" flashed")

			// Check if we need to flash more
			if msg.Step == firmware.FlashLeft || msg.Step == firmware.FlashResetLeft {
				// Wait for device reconnect for right half
				m.flashTarget = "RIGHT half"
				m.state = StateWaitingDevice
				m.flashStep = firmware.FlashRight
				if msg.Step == firmware.FlashResetLeft {
					m.flashStep = firmware.FlashResetRight
				}
				return m, m.watcher.WaitForDevice()
			}

			// All done
			m.state = StateComplete
			m.logPanel.Add(LogSuccess, "Flash complete")
		} else {
			m.logPanel.Add(LogError, "Flash failed: "+msg.Error.Error())
			m.state = StateIdle
		}
		return m, nil

	case tickMsg:
		// Refresh spinner
		return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return tickMsg{}
		})
	}

	return m, nil
}

type tickMsg struct{}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keys - ctrl+c always quits
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "q":
		if !m.showDialog && (m.state == StateIdle || m.state == StateComplete) {
			return m, tea.Quit
		}
	case "?":
		if m.state == StateIdle {
			m.showHelp = !m.showHelp
		}
		return m, nil
	case "esc":
		if m.showHelp {
			m.showHelp = false
			return m, nil
		}
		if m.showDialog {
			m.showDialog = false
			m.confirmDialog = nil
			return m, nil
		}
		if m.state == StateWaitingDevice {
			m.watcher.Cancel()
			m.state = StateIdle
			m.logPanel.Add(LogInfo, "Cancelled")
			return m, nil
		}
		if m.state == StateComplete {
			m.state = StateIdle
			m.completedSteps = nil
			return m, nil
		}
	}

	// Dialog keys
	if m.showDialog && m.confirmDialog != nil {
		switch msg.String() {
		case "left", "h":
			m.confirmDialog.MoveLeft()
		case "right", "l":
			m.confirmDialog.MoveRight()
		case "enter":
			if m.confirmDialog.Selected() == DialogConfirm {
				m.showDialog = false
				return m.startFactoryReset()
			}
			m.showDialog = false
			m.confirmDialog = nil
		}
		return m, nil
	}

	// Help overlay blocks other keys
	if m.showHelp {
		return m, nil
	}

	// State-specific keys
	switch m.state {
	case StateIdle:
		return m.handleIdleKey(msg)
	case StateComplete:
		if msg.String() == "enter" {
			m.state = StateIdle
			m.completedSteps = nil
		}
	}

	return m, nil
}

func (m *Model) handleIdleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	// Navigation
	case "up", "k":
		if m.activePanel == PanelFirmware {
			m.firmwarePanel.MoveUp()
		}
	case "down", "j":
		if m.activePanel == PanelFirmware {
			m.firmwarePanel.MoveDown()
		}
	case "tab":
		m.activePanel = (m.activePanel + 1) % 3
	case "1":
		m.activePanel = PanelFirmware
	case "2":
		m.activePanel = PanelStatus
	case "3":
		m.activePanel = PanelLog

	// Actions
	case "b":
		return m.startBuild()
	case "f", "enter":
		if m.firmwarePanel.Selected() != nil {
			return m.prepareFlash()
		}
	case "r":
		build := m.firmwarePanel.Selected()
		if build != nil && build.HasReset {
			m.confirmDialog = FactoryResetDialog()
			m.confirmDialog.SetSize(m.width, m.height)
			m.showDialog = true
		}
	}

	return m, nil
}

func (m *Model) startBuild() (tea.Model, tea.Cmd) {
	m.state = StateBuilding
	m.buildStep = firmware.StepActivateEnv
	m.buildPercent = 0
	m.startTime = time.Now()
	m.logPanel.Add(LogInfo, "Build started")

	return m, tea.Batch(
		m.builder.Build("both"),
		tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return tickMsg{}
		}),
	)
}

func (m *Model) prepareFlash() (tea.Model, tea.Cmd) {
	build := m.firmwarePanel.Selected()
	if build == nil || !build.HasLeft || !build.HasRight {
		m.logPanel.Add(LogError, "Missing firmware files")
		return m, nil
	}

	m.completedSteps = nil
	m.flashTarget = "LEFT half"
	m.flashStep = firmware.FlashLeft
	m.startTime = time.Now()

	if m.deviceStatus == device.Connected {
		return m.startFlash()
	}

	m.state = StateWaitingDevice
	m.logPanel.Add(LogInfo, "Waiting for device...")
	return m, tea.Batch(
		m.watcher.WaitForDevice(),
		tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return tickMsg{}
		}),
	)
}

func (m *Model) startFlash() (tea.Model, tea.Cmd) {
	build := m.firmwarePanel.Selected()
	if build == nil {
		return m, nil
	}

	m.state = StateFlashing
	m.flashPercent = 0
	m.logPanel.Add(LogInfo, "Flashing "+m.flashStep.String())

	var filePath string
	switch m.flashStep {
	case firmware.FlashLeft:
		filePath = filepath.Join(build.Path, "corne_left.uf2")
	case firmware.FlashRight:
		filePath = filepath.Join(build.Path, "corne_right.uf2")
	case firmware.FlashResetLeft, firmware.FlashResetRight:
		filePath = filepath.Join(build.Path, "settings_reset.uf2")
	}

	return m, tea.Batch(
		m.flasher.Flash(m.flashStep, filePath),
		tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return tickMsg{}
		}),
	)
}

func (m *Model) startFactoryReset() (tea.Model, tea.Cmd) {
	build := m.firmwarePanel.Selected()
	if build == nil || !build.HasReset {
		return m, nil
	}

	m.completedSteps = nil
	m.flashTarget = "LEFT half (reset)"
	m.flashStep = firmware.FlashResetLeft
	m.startTime = time.Now()
	m.logPanel.Add(LogWarning, "Factory reset started")

	if m.deviceStatus == device.Connected {
		m.state = StateFlashing
		resetPath := filepath.Join(build.Path, "settings_reset.uf2")
		return m, m.flasher.Flash(firmware.FlashResetLeft, resetPath)
	}

	m.state = StateWaitingDevice
	return m, m.watcher.WaitForDevice()
}

func (m *Model) updatePanelSizes() {
	// Account for header (2), footer (1), and borders
	contentHeight := m.height - 4

	// Panel widths (30% / 40% / 30%)
	leftWidth := m.width * 30 / 100
	centerWidth := m.width * 40 / 100
	rightWidth := m.width - leftWidth - centerWidth - 6 // Account for borders

	m.firmwarePanel.SetSize(leftWidth, contentHeight)
	m.statusPanel.SetSize(centerWidth, contentHeight)
	m.logPanel.SetSize(rightWidth, contentHeight)
	m.helpOverlay.SetSize(m.width, m.height)
	if m.confirmDialog != nil {
		m.confirmDialog.SetSize(m.width, m.height)
	}
}

// View renders the UI
func (m *Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	var s strings.Builder

	// Header
	s.WriteString(m.renderHeader())
	s.WriteString("\n")

	// Main content
	s.WriteString(m.renderPanels())

	// Footer
	s.WriteString("\n")
	s.WriteString(m.renderFooter())

	// Overlay
	if m.showHelp {
		return m.helpOverlay.View()
	}
	if m.showDialog && m.confirmDialog != nil {
		return m.confirmDialog.View()
	}

	return s.String()
}

func (m *Model) renderHeader() string {
	title := TitleStyle.Render("⌨  CORNE FLASH UTILITY")

	// Device status
	var statusIcon, statusText string
	switch m.deviceStatus {
	case device.Connected:
		statusIcon = SuccessStyle.Render(StatusConnected)
		statusText = "NICENANO Connected"
	case device.Waiting:
		statusIcon = WarningStyle.Render(StatusWaiting)
		statusText = "NICENANO Waiting..."
	default:
		statusIcon = DimStyle.Render(StatusDisconnected)
		statusText = "NICENANO Disconnected"
	}
	status := statusIcon + " " + statusText

	version := DimStyle.Render("v1.0.0")

	// Calculate spacing
	leftPart := title
	rightPart := status + "   " + version
	spacing := m.width - lipgloss.Width(leftPart) - lipgloss.Width(rightPart) - 2
	if spacing < 1 {
		spacing = 1
	}

	headerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Width(m.width - 2)

	content := leftPart + strings.Repeat(" ", spacing) + rightPart
	return headerStyle.Render(content)
}

func (m *Model) renderPanels() string {
	// Panel widths
	leftWidth := m.width * 30 / 100
	centerWidth := m.width * 40 / 100
	rightWidth := m.width - leftWidth - centerWidth - 6

	// Content height
	contentHeight := m.height - 6

	// Firmware panel
	firmwareStyle := PanelStyle.Width(leftWidth).Height(contentHeight)
	if m.activePanel == PanelFirmware {
		firmwareStyle = ActivePanelStyle.Width(leftWidth).Height(contentHeight)
	}
	firmwareTitle := " Firmware "
	firmwareContent := m.firmwarePanel.View()
	firmwarePanel := firmwareStyle.Render(AccentStyle.Render(firmwareTitle) + "\n\n" + firmwareContent)

	// Status panel
	statusStyle := PanelStyle.Width(centerWidth).Height(contentHeight)
	if m.activePanel == PanelStatus {
		statusStyle = ActivePanelStyle.Width(centerWidth).Height(contentHeight)
	}
	statusTitle := " Status "
	var statusContent string
	switch m.state {
	case StateIdle:
		statusContent = m.statusPanel.ViewIdle(m.firmwarePanel.Selected())
	case StateBuilding:
		statusContent = m.statusPanel.ViewBuilding(m.buildStep, m.buildPercent)
	case StateWaitingDevice:
		statusContent = m.statusPanel.ViewWaiting(m.flashTarget)
	case StateFlashing:
		build := m.firmwarePanel.Selected()
		filename := ""
		if build != nil {
			switch m.flashStep {
			case firmware.FlashLeft:
				filename = "corne_left.uf2"
			case firmware.FlashRight:
				filename = "corne_right.uf2"
			default:
				filename = "settings_reset.uf2"
			}
		}
		statusContent = m.statusPanel.ViewFlashing(m.flashStep, m.flashPercent, filename)
	case StateComplete:
		duration := time.Since(m.startTime)
		statusContent = m.statusPanel.ViewComplete(duration, m.completedSteps)
	}
	statusPanel := statusStyle.Render(AccentStyle.Render(statusTitle) + "\n\n" + statusContent)

	// Log panel
	logStyle := PanelStyle.Width(rightWidth).Height(contentHeight)
	if m.activePanel == PanelLog {
		logStyle = ActivePanelStyle.Width(rightWidth).Height(contentHeight)
	}
	logTitle := " Log "
	logContent := m.logPanel.View()
	logPanel := logStyle.Render(AccentStyle.Render(logTitle) + "\n\n" + logContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, firmwarePanel, statusPanel, logPanel)
}

func (m *Model) renderFooter() string {
	var hints []string

	switch m.state {
	case StateIdle:
		hints = []string{
			"↑/↓ Navigate",
			"Enter Select",
			"b Build",
			"f Flash",
			"r Factory Reset",
			"q Quit",
		}
	case StateBuilding:
		hints = []string{"Building..."}
	case StateWaitingDevice:
		hints = []string{"Waiting for device... Double-tap reset button", "Esc Cancel"}
	case StateFlashing:
		hints = []string{"Flashing... Do not disconnect device"}
	case StateComplete:
		hints = []string{"Enter Continue", "q Quit"}
	}

	left := DimStyle.Render(strings.Join(hints, "   "))
	right := DimStyle.Render("? Help")

	spacing := m.width - lipgloss.Width(left) - lipgloss.Width(right) - 2
	if spacing < 1 {
		spacing = 1
	}

	return " " + left + strings.Repeat(" ", spacing) + right
}

func formatInt(n int) string {
	if n == 0 {
		return "0"
	}
	var result []byte
	for n > 0 {
		result = append([]byte{byte('0' + n%10)}, result...)
		n /= 10
	}
	return string(result)
}
