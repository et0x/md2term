package theme

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

type Theme struct {
	Name          string `yaml:"name"`
	Heading1      string `yaml:"heading1"`
	Heading2      string `yaml:"heading2"`
	Heading3      string `yaml:"heading3"`
	Heading4      string `yaml:"heading4"`
	Heading5      string `yaml:"heading5"`
	Heading6      string `yaml:"heading6"`
	ListItem      string `yaml:"listItem"`
	UncheckedItem string `yaml:"uncheckedItem"`
	CheckedItem   string `yaml:"checkedItem"`
	TableBorder   string `yaml:"tableBorder"`
	TableCell     string `yaml:"tableCell"`
	TableHeader   struct {
		Fg string `yaml:"fg"`
		Bg string `yaml:"bg"`
	} `yaml:"tableHeader"`
	EvenRow   string `yaml:"evenRow"`
	OddRow    string `yaml:"oddRow"`
	CodeBlock struct {
		Fg string `yaml:"fg"`
		Bg string `yaml:"bg"`
	} `yaml:"codeBlock"`
	Italic     string `yaml:"italic"`
	Bold       string `yaml:"bold"`
	Link       string `yaml:"link"`
	Image      string `yaml:"image"`
	Blockquote string `yaml:"blockquote"`
	InlineCode struct {
		Fg string `yaml:"fg"`
		Bg string `yaml:"bg"`
	} `yaml:"inlineCode"`
}

type Manager struct {
	CurrentTheme Theme
	Styles       Styles
}

type Styles struct {
	Heading1      lipgloss.Style
	Heading2      lipgloss.Style
	Heading3      lipgloss.Style
	Heading4      lipgloss.Style
	Heading5      lipgloss.Style
	Heading6      lipgloss.Style
	ListItem      lipgloss.Style
	UncheckedItem lipgloss.Style
	CheckedItem   lipgloss.Style
	TableBorder   lipgloss.Style
	TableCell     lipgloss.Style
	TableHeader   lipgloss.Style
	EvenRow       lipgloss.Style
	OddRow        lipgloss.Style
	CodeBlock     lipgloss.Style
	Italic        lipgloss.Style
	Bold          lipgloss.Style
	Link          lipgloss.Style
	Image         lipgloss.Style
	Blockquote    lipgloss.Style
	InlineCode    lipgloss.Style
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) LoadTheme(themeName string) error {
	// Read the themes.yaml file
	data, err := os.ReadFile("themes.yaml")
	if err != nil {
		return fmt.Errorf("error reading theme file: %w", err)
	}

	var themes map[string]Theme
	err = yaml.Unmarshal(data, &themes)
	if err != nil {
		return fmt.Errorf("error parsing theme file: %w", err)
	}

	theme, ok := themes[themeName]
	if !ok {
		return fmt.Errorf("theme '%s' not found", themeName)
	}

	m.CurrentTheme = theme
	m.applyTheme()
	return nil
}

func (m *Manager) applyTheme() {
	m.Styles = Styles{
		Heading1:      lipgloss.NewStyle().Foreground(lipgloss.Color(m.CurrentTheme.Heading1)).Bold(true).Underline(true),
		Heading2:      lipgloss.NewStyle().Foreground(lipgloss.Color(m.CurrentTheme.Heading2)).Bold(true),
		Heading3:      lipgloss.NewStyle().Foreground(lipgloss.Color(m.CurrentTheme.Heading3)).Bold(true),
		Heading4:      lipgloss.NewStyle().Foreground(lipgloss.Color(m.CurrentTheme.Heading4)).Bold(true),
		Heading5:      lipgloss.NewStyle().Foreground(lipgloss.Color(m.CurrentTheme.Heading5)).Bold(true),
		Heading6:      lipgloss.NewStyle().Foreground(lipgloss.Color(m.CurrentTheme.Heading6)).Bold(true),
		ListItem:      lipgloss.NewStyle().Foreground(lipgloss.Color(m.CurrentTheme.ListItem)),
		UncheckedItem: lipgloss.NewStyle().Foreground(lipgloss.Color(m.CurrentTheme.UncheckedItem)),
		CheckedItem:   lipgloss.NewStyle().Foreground(lipgloss.Color(m.CurrentTheme.CheckedItem)),
		TableBorder:   lipgloss.NewStyle().Foreground(lipgloss.Color(m.CurrentTheme.TableBorder)),
		TableCell:     lipgloss.NewStyle().Foreground(lipgloss.Color(m.CurrentTheme.TableCell)).PaddingLeft(1).PaddingRight(1),
		TableHeader:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(m.CurrentTheme.TableHeader.Fg)).Background(lipgloss.Color(m.CurrentTheme.TableHeader.Bg)),
		EvenRow:       lipgloss.NewStyle().Foreground(lipgloss.Color(m.CurrentTheme.EvenRow)),
		OddRow:        lipgloss.NewStyle().Foreground(lipgloss.Color(m.CurrentTheme.OddRow)),
		CodeBlock:     lipgloss.NewStyle().Foreground(lipgloss.Color(m.CurrentTheme.CodeBlock.Fg)).Background(lipgloss.Color(m.CurrentTheme.CodeBlock.Bg)).Padding(1),
		Italic:        lipgloss.NewStyle().Italic(true),
		Bold:          lipgloss.NewStyle().Bold(true),
		Link:          lipgloss.NewStyle().Foreground(lipgloss.Color(m.CurrentTheme.Link)).Underline(true),
		Image:         lipgloss.NewStyle().Foreground(lipgloss.Color(m.CurrentTheme.Image)),
		Blockquote:    lipgloss.NewStyle().Foreground(lipgloss.Color(m.CurrentTheme.Blockquote)).PaddingLeft(2),
		InlineCode:    lipgloss.NewStyle().Foreground(lipgloss.Color(m.CurrentTheme.InlineCode.Fg)).Background(lipgloss.Color(m.CurrentTheme.InlineCode.Bg)),
	}
}

func (m *Manager) SaveCurrentTheme() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("error getting user config directory: %w", err)
	}

	md2termDir := filepath.Join(configDir, "md2term")
	err = os.MkdirAll(md2termDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	themeFile := filepath.Join(md2termDir, "current_theme")
	err = os.WriteFile(themeFile, []byte(m.CurrentTheme.Name), 0644)
	if err != nil {
		return fmt.Errorf("error saving current theme: %w", err)
	}

	return nil
}

func (m *Manager) LoadSavedTheme() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("error getting user config directory: %w", err)
	}

	themeFile := filepath.Join(configDir, "md2term", "current_theme")
	data, err := os.ReadFile(themeFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "default", nil
		}
		return "", fmt.Errorf("error reading saved theme: %w", err)
	}

	return string(data), nil
}
