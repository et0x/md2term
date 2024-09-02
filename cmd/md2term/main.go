package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	heading1Style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			Underline(true)

	heading2Style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4500")).
			Bold(true)

	heading3Style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Bold(true)

	listItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF"))

	uncheckedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#1E90FF"))

	checkedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00FF00"))

	tableBorderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#808080"))

	tableCellStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF")).
			PaddingLeft(1).
			PaddingRight(1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#333333"))

	evenRowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF"))

	oddRowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00CCCC"))

	codeBlockStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#333333")).
			Padding(1)
)

func main() {
	var input io.Reader

	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Println("Error opening file:", err)
			os.Exit(1)
		}
		defer file.Close()
		input = file
	} else {
		input = os.Stdin
	}

	var tableLines []string
	inTable := false
	inCodeBlock := false
	var codeBlockLines []string

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()

		if inCodeBlock {
			if strings.HasPrefix(line, "```") {
				formatCodeBlock(codeBlockLines)
				codeBlockLines = nil
				inCodeBlock = false
			} else {
				codeBlockLines = append(codeBlockLines, line)
			}
			continue
		}

		if strings.HasPrefix(line, "```") {
			inCodeBlock = true
			continue
		}

		if strings.HasPrefix(line, "|") {
			if !inTable {
				inTable = true
			}
			tableLines = append(tableLines, line)
		} else {
			if inTable {
				formatTable(tableLines)
				tableLines = nil
				inTable = false
			}
			formatLine(line)
		}
	}

	if inTable {
		formatTable(tableLines)
	}

	if inCodeBlock {
		formatCodeBlock(codeBlockLines)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
}

func formatLine(line string) {
	switch {
	case strings.HasPrefix(line, "# "):
		fmt.Println(heading1Style.Render(line[2:]))
	case strings.HasPrefix(line, "## "):
		fmt.Println(heading2Style.Render(line[3:]))
	case strings.HasPrefix(line, "### "):
		fmt.Println(heading3Style.Render(line[4:]))
	case strings.HasPrefix(line, "- "):
		fmt.Println("  •", listItemStyle.Render(line[2:]))
	case strings.HasPrefix(line, "- [ ] "):
		fmt.Println("  ☐", uncheckedItemStyle.Render(line[6:]))
	case strings.HasPrefix(line, "- [x] "):
		fmt.Println("  ☑", checkedItemStyle.Render(line[6:]))
	case strings.HasPrefix(line, "|"):
		formatTable([]string{line})
	default:
		fmt.Println(line)
	}
}

func formatTable(lines []string) {
	var rows [][]string
	for i, line := range lines {
		if i == 1 && strings.Contains(line, "---") {
			continue // Skip the separator line
		}
		cells := strings.Split(strings.Trim(line, "|"), "|")
		for i, cell := range cells {
			cells[i] = strings.TrimSpace(cell)
		}
		rows = append(rows, cells)
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return headerStyle
			case row%2 == 0:
				return evenRowStyle
			default:
				return oddRowStyle
			}
		})

	if len(rows) > 0 {
		t.Headers(rows[0]...)
		t.Rows(rows[1:]...)
	}

	fmt.Println(t.Render())
}

func formatCodeBlock(lines []string) {
	codeBlock := strings.Join(lines, "\n")
	fmt.Println(codeBlockStyle.Render(codeBlock))
}
