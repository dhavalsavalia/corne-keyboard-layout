package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// DialogOption represents a dialog button
type DialogOption int

const (
	DialogConfirm DialogOption = iota
	DialogCancel
)

// ConfirmDialog renders a confirmation dialog
type ConfirmDialog struct {
	title    string
	message  []string
	selected DialogOption
	width    int
	height   int
}

// NewConfirmDialog creates a new confirmation dialog
func NewConfirmDialog(title string, message []string) *ConfirmDialog {
	return &ConfirmDialog{
		title:    title,
		message:  message,
		selected: DialogCancel, // Default to cancel for safety
	}
}

// SetSize sets dialog dimensions
func (d *ConfirmDialog) SetSize(width, height int) {
	d.width = width
	d.height = height
}

// MoveLeft moves selection left (to confirm)
func (d *ConfirmDialog) MoveLeft() {
	d.selected = DialogConfirm
}

// MoveRight moves selection right (to cancel)
func (d *ConfirmDialog) MoveRight() {
	d.selected = DialogCancel
}

// Selected returns the selected option
func (d *ConfirmDialog) Selected() DialogOption {
	return d.selected
}

// View renders the dialog
func (d *ConfirmDialog) View() string {
	var lines []string

	// Title with warning icon
	title := WarningStyle.Render("⚠  " + d.title)
	lines = append(lines, title)
	lines = append(lines, "")

	// Message lines
	for _, msg := range d.message {
		lines = append(lines, msg)
	}
	lines = append(lines, "")

	// Buttons
	confirmStyle := lipgloss.NewStyle().Padding(0, 2)
	cancelStyle := lipgloss.NewStyle().Padding(0, 2)

	if d.selected == DialogConfirm {
		confirmStyle = confirmStyle.Background(ColorPurple).Foreground(lipgloss.Color("0"))
	}
	if d.selected == DialogCancel {
		cancelStyle = cancelStyle.Background(ColorPurple).Foreground(lipgloss.Color("0"))
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Center,
		confirmStyle.Render("Yes, proceed"),
		"  ",
		cancelStyle.Render("Cancel"),
	)
	lines = append(lines, buttons)

	content := strings.Join(lines, "\n")

	// Box style
	boxWidth := 40
	if boxWidth > d.width-10 {
		boxWidth = d.width - 10
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorYellow).
		Padding(1, 2).
		Width(boxWidth)

	box := boxStyle.Render(content)

	// Center the box
	boxHeight := lipgloss.Height(box)
	topPadding := (d.height - boxHeight) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	leftPadding := (d.width - boxWidth - 4) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}

	var result []string
	for i := 0; i < topPadding; i++ {
		result = append(result, "")
	}

	for _, line := range strings.Split(box, "\n") {
		result = append(result, strings.Repeat(" ", leftPadding)+line)
	}

	return strings.Join(result, "\n")
}

// FactoryResetDialog creates the factory reset confirmation dialog
func FactoryResetDialog() *ConfirmDialog {
	return NewConfirmDialog("FACTORY RESET", []string{
		"This will:",
		"• Clear all Bluetooth bonds",
		"• Reset keyboard settings",
		"• Require re-pairing",
		"",
		"Have you unpaired from all",
		"Bluetooth devices?",
	})
}
