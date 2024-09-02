package formatter

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/et0x/md2term/internal/theme"
)

type Formatter struct {
	themeManager *theme.Manager
}

func NewFormatter(tm *theme.Manager) *Formatter {
	return &Formatter{themeManager: tm}
}

func (f *Formatter) ProcessMarkdown(input io.Reader) error {
	var tableLines []string
	inTable := false
	inCodeBlock := false
	var codeBlockLines []string
	var blockquoteLines []string
	inBlockquote := false

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case inCodeBlock:
			if strings.HasPrefix(line, "```") {
				f.formatCodeBlock(codeBlockLines)
				codeBlockLines = nil
				inCodeBlock = false
			} else {
				codeBlockLines = append(codeBlockLines, line)
			}
		case inBlockquote:
			if strings.HasPrefix(line, ">") || line == "" {
				blockquoteLines = append(blockquoteLines, line)
			} else {
				f.formatBlockquote(blockquoteLines)
				blockquoteLines = nil
				inBlockquote = false
				f.formatLine(line)
			}
		case strings.HasPrefix(line, "```"):
			inCodeBlock = true
		case strings.HasPrefix(line, ">"):
			inBlockquote = true
			blockquoteLines = append(blockquoteLines, line)
		case strings.HasPrefix(line, "|"):
			if !inTable {
				inTable = true
			}
			tableLines = append(tableLines, line)
		default:
			if inTable {
				f.formatTable(tableLines)
				tableLines = nil
				inTable = false
			}
			f.formatLine(line)
		}
	}

	// Handle any remaining elements
	if inTable {
		f.formatTable(tableLines)
	}
	if inCodeBlock {
		f.formatCodeBlock(codeBlockLines)
	}
	if inBlockquote {
		f.formatBlockquote(blockquoteLines)
	}

	return scanner.Err()
}

func (f *Formatter) formatLine(line string) {
	switch {
	case strings.HasPrefix(line, "# "):
		fmt.Println(f.themeManager.Styles.Heading1.Render(line[2:]))
	case strings.HasPrefix(line, "## "):
		fmt.Println(f.themeManager.Styles.Heading2.Render(line[3:]))
	case strings.HasPrefix(line, "### "):
		fmt.Println(f.themeManager.Styles.Heading3.Render(line[4:]))
	case strings.HasPrefix(line, "#### "):
		fmt.Println(f.themeManager.Styles.Heading4.Render(line[5:]))
	case strings.HasPrefix(line, "##### "):
		fmt.Println(f.themeManager.Styles.Heading5.Render(line[6:]))
	case strings.HasPrefix(line, "###### "):
		fmt.Println(f.themeManager.Styles.Heading6.Render(line[7:]))
	case strings.HasPrefix(line, "- "):
		fmt.Println("  •", f.themeManager.Styles.ListItem.Render(line[2:]))
	case strings.HasPrefix(line, "- [ ] "):
		fmt.Println("  ☐", f.themeManager.Styles.UncheckedItem.Render(line[6:]))
	case strings.HasPrefix(line, "- [x] "):
		fmt.Println("  ☑", f.themeManager.Styles.CheckedItem.Render(line[6:]))
	case strings.HasPrefix(line, "|"):
		f.formatTable([]string{line})
	case strings.HasPrefix(line, ">"):
		f.formatBlockquote([]string{line})
	case strings.HasPrefix(line, "!["):
		fmt.Println(f.formatImage(line))
	default:
		fmt.Println(f.formatInline(line))
	}
}

func (f *Formatter) formatInline(line string) string {
	// Handle bold and italic
	line = boldItalicRegex.ReplaceAllStringFunc(line, func(s string) string {
		return f.themeManager.Styles.Bold.Render(f.themeManager.Styles.Italic.Render(s[3 : len(s)-3]))
	})
	line = boldRegex.ReplaceAllStringFunc(line, func(s string) string {
		return f.themeManager.Styles.Bold.Render(s[2 : len(s)-2])
	})
	line = italicRegex.ReplaceAllStringFunc(line, func(s string) string {
		return f.themeManager.Styles.Italic.Render(s[1 : len(s)-1])
	})

	// Handle links
	line = linkRegex.ReplaceAllStringFunc(line, func(s string) string {
		parts := linkRegex.FindStringSubmatch(s)
		return fmt.Sprintf("%s (%s)", f.themeManager.Styles.Link.Render(parts[1]), parts[2])
	})

	// Handle inline code
	line = inlineCodeRegex.ReplaceAllStringFunc(line, func(s string) string {
		return f.themeManager.Styles.InlineCode.Render(s[1 : len(s)-1])
	})

	return line
}

func (f *Formatter) formatBlockquote(lines []string) {
	for _, line := range lines {
		indent := 0
		for strings.HasPrefix(line, ">") {
			indent++
			line = strings.TrimPrefix(line, ">")
			line = strings.TrimSpace(line)
		}
		if line == "" {
			fmt.Println(f.themeManager.Styles.Blockquote.Copy().PaddingLeft(indent * 2).Render(""))
		} else {
			fmt.Println(f.themeManager.Styles.Blockquote.Copy().PaddingLeft(indent * 2).Render(line))
		}
	}
}

func (f *Formatter) formatImage(line string) string {
	parts := imageRegex.FindStringSubmatch(line)
	if len(parts) == 4 {
		return f.themeManager.Styles.Image.Render(fmt.Sprintf("[Image: %s (%s)]", parts[1], parts[2]))
	}
	return line
}

func (f *Formatter) formatTable(lines []string) {
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
		BorderStyle(f.themeManager.Styles.TableBorder).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return f.themeManager.Styles.TableHeader
			case row%2 == 0:
				return f.themeManager.Styles.EvenRow
			default:
				return f.themeManager.Styles.OddRow
			}
		})

	if len(rows) > 0 {
		t.Headers(rows[0]...)
		t.Rows(rows[1:]...)
	}

	fmt.Println(t.Render())
}

func (f *Formatter) formatCodeBlock(lines []string) {
	codeBlock := strings.Join(lines, "\n")
	fmt.Println(f.themeManager.Styles.CodeBlock.Render(codeBlock))
}

var (
	boldItalicRegex = regexp.MustCompile(`(\*\*\*|___)(.+?)(\*\*\*|___)`)
	boldRegex       = regexp.MustCompile(`(\*\*|__)(.+?)(\*\*|__)`)
	italicRegex     = regexp.MustCompile(`(\*|_)(.+?)(\*|_)`)
	linkRegex       = regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`)
	imageRegex      = regexp.MustCompile(`!\[([^\]]*)\]\(([^\s]+)\s*"?([^"]*)"?\)`)
	inlineCodeRegex = regexp.MustCompile("`([^`]+)`")
)
