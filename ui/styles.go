package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Color palette - retro green terminal style
var (
	// Base colors
	ColorPrimary   = lipgloss.Color("#00FF00") // Bright green
	ColorSecondary = lipgloss.Color("#00AA00") // Darker green
	ColorAccent    = lipgloss.Color("#FFFF00") // Yellow for highlights
	ColorDanger    = lipgloss.Color("#FF0000") // Red for warnings
	ColorInfo      = lipgloss.Color("#00FFFF") // Cyan for info
	ColorMuted     = lipgloss.Color("#888888") // Gray for muted text

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2).
			Margin(1, 0)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent).
			Padding(0, 1)

	MenuItemStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Padding(0, 2)

	SelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(ColorPrimary).
			Padding(0, 2)

	InfoStyle = lipgloss.NewStyle().
			Foreground(ColorInfo).
			Padding(0, 1)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			Padding(0, 1)

	DangerStyle = lipgloss.NewStyle().
			Foreground(ColorDanger).
			Bold(true).
			Padding(0, 1)

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

	MutedStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Padding(0, 1)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorSecondary).
			Padding(1, 2).
			Margin(1, 0)

	StatusBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorPrimary).
			Padding(0, 1).
			Margin(0, 1)

	ButtonStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(ColorPrimary).
			Padding(0, 3).
			Margin(0, 1)

	InputStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Background(lipgloss.Color("#1a1a1a")).
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorSecondary)

	// ASCII art styles
	AsciiStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)
)

// Border characters for box drawing
const (
	BorderTL = "╔"
	BorderTR = "╗"
	BorderBL = "╚"
	BorderBR = "╝"
	BorderH  = "═"
	BorderV  = "║"
)

// DrawBox creates a bordered box with content
func DrawBox(width int, title string, content string) string {
	var result string

	// Top border with title
	result += BorderTL
	for i := 0; i < width-2; i++ {
		result += BorderH
	}
	result += BorderTR + "\n"

	// Title line
	if title != "" {
		titleLine := fmt.Sprintf("  %s  ", title)
		padding := (width - len(titleLine) - 2) / 2
		result += BorderV
		for i := 0; i < padding; i++ {
			result += " "
		}
		result += titleLine
		for i := 0; i < width-padding-len(titleLine)-2; i++ {
			result += " "
		}
		result += BorderV + "\n"
	}

	// Content lines
	for _, line := range strings.Split(content, "\n") {
		result += BorderV
		result += " "
		result += line
		// Pad to width
		for i := len(line) + 1; i < width-2; i++ {
			result += " "
		}
		result += BorderV + "\n"
	}

	// Bottom border
	result += BorderBL
	for i := 0; i < width-2; i++ {
		result += BorderH
	}
	result += BorderBR

	return result
}

// CenterText centers text within a given width
func CenterText(text string, width int) string {
	padding := (width - len(text)) / 2
	if padding < 0 {
		padding = 0
	}
	result := ""
	for i := 0; i < padding; i++ {
		result += " "
	}
	result += text
	return result
}

// ProgressBar creates a text-based progress bar
func ProgressBar(current, max, width int) string {
	if max == 0 {
		max = 1
	}
	percent := float64(current) / float64(max)
	filled := int(percent * float64(width))

	bar := "["
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	bar += "]"
	return bar
}
