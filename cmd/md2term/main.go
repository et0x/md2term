package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
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

	italicStyle = lipgloss.NewStyle().
			Italic(true)

	boldStyle = lipgloss.NewStyle().
			Bold(true)

	linkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0000FF")).
			Underline(true)

	imageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#008000"))

	blockquoteStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080")).
			PaddingLeft(2)

	inlineCodeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF69B4")).
			Background(lipgloss.Color("#F0F0F0"))

	heading4Style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true)

	heading5Style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DAA520")).
			Bold(true)

	heading6Style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#B8860B")).
			Bold(true)
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
	var blockquoteLines []string
	inBlockquote := false

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
		} else if inBlockquote {
			if strings.HasPrefix(line, ">") || line == "" {
				blockquoteLines = append(blockquoteLines, line)
			} else {
				formatBlockquote(blockquoteLines)
				blockquoteLines = nil
				inBlockquote = false
				formatLine(line)
			}
		} else if strings.HasPrefix(line, "```") {
			inCodeBlock = true
		} else if strings.HasPrefix(line, ">") {
			inBlockquote = true
			blockquoteLines = append(blockquoteLines, line)
		} else if strings.HasPrefix(line, "|") {
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

	// Handle any remaining elements
	if inTable {
		formatTable(tableLines)
	}
	if inCodeBlock {
		formatCodeBlock(codeBlockLines)
	}
	if inBlockquote {
		formatBlockquote(blockquoteLines)
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
	case strings.HasPrefix(line, "#### "):
		fmt.Println(heading4Style.Render(line[5:]))
	case strings.HasPrefix(line, "##### "):
		fmt.Println(heading5Style.Render(line[6:]))
	case strings.HasPrefix(line, "###### "):
		fmt.Println(heading6Style.Render(line[7:]))
	case strings.HasPrefix(line, "- "):
		fmt.Println("  •", listItemStyle.Render(line[2:]))
	case strings.HasPrefix(line, "- [ ] "):
		fmt.Println("  ☐", uncheckedItemStyle.Render(line[6:]))
	case strings.HasPrefix(line, "- [x] "):
		fmt.Println("  ☑", checkedItemStyle.Render(line[6:]))
	case strings.HasPrefix(line, "|"):
		formatTable([]string{line})
	case strings.HasPrefix(line, ">"):
		formatBlockquote([]string{line})
	case strings.HasPrefix(line, "!["):
		fmt.Println(formatImage(line))
	default:
		fmt.Println(formatInline(line))
	}
}

func formatInline(line string) string {
	// Handle bold and italic
	line = boldItalicRegex.ReplaceAllStringFunc(line, func(s string) string {
		return boldStyle.Render(italicStyle.Render(s[3 : len(s)-3]))
	})
	line = boldRegex.ReplaceAllStringFunc(line, func(s string) string {
		return boldStyle.Render(s[2 : len(s)-2])
	})
	line = italicRegex.ReplaceAllStringFunc(line, func(s string) string {
		return italicStyle.Render(s[1 : len(s)-1])
	})

	// Handle links
	line = linkRegex.ReplaceAllStringFunc(line, func(s string) string {
		parts := linkRegex.FindStringSubmatch(s)
		return fmt.Sprintf("%s (%s)", linkStyle.Render(parts[1]), parts[2])
	})

	// Handle inline code
	line = inlineCodeRegex.ReplaceAllStringFunc(line, func(s string) string {
		return inlineCodeStyle.Render(s[1 : len(s)-1])
	})

	return line
}

func formatBlockquote(lines []string) {
	for _, line := range lines {
		indent := 0
		for strings.HasPrefix(line, ">") {
			indent++
			line = strings.TrimPrefix(line, ">")
			line = strings.TrimSpace(line)
		}
		if line == "" {
			fmt.Println(blockquoteStyle.Copy().PaddingLeft(indent * 2).Render(""))
		} else {
			fmt.Println(blockquoteStyle.Copy().PaddingLeft(indent * 2).Render(line))
		}
	}
}

func formatImage(line string) string {
	parts := imageRegex.FindStringSubmatch(line)
	if len(parts) == 4 {
		return imageStyle.Render(fmt.Sprintf("[Image: %s (%s)]", parts[1], parts[2]))
	}
	return line
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

var (
	boldItalicRegex = regexp.MustCompile(`(\*\*\*|___)(.+?)(\*\*\*|___)`)
	boldRegex       = regexp.MustCompile(`(\*\*|__)(.+?)(\*\*|__)`)
	italicRegex     = regexp.MustCompile(`(\*|_)(.+?)(\*|_)`)
	linkRegex       = regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`)
	imageRegex      = regexp.MustCompile(`!\[([^\]]*)\]\(([^\s]+)\s*"?([^"]*)"?\)`)
	inlineCodeRegex = regexp.MustCompile("`([^`]+)`")
)
