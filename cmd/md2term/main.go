package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/et0x/md2term/internal/formatter"
	"github.com/et0x/md2term/internal/theme"
)

// Version is the current version of md2term
const Version = "1.0.1"

func main() {
	// Parse command-line flags
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	themeFlag := flag.String("theme", "", "Specify the theme to use (default, dark, or light)")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("md2term version %s\n", Version)
		return
	}

	themeManager := theme.NewManager()

	// Load the theme
	themeName := *themeFlag
	if themeName == "" {
		var err error
		themeName, err = themeManager.LoadSavedTheme()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading saved theme: %v\n", err)
			os.Exit(1)
		}
	}

	err := themeManager.LoadTheme(themeName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading theme: %v\n", err)
		os.Exit(1)
	}

	// Save the current theme if specified
	if *themeFlag != "" {
		err = themeManager.SaveCurrentTheme()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving current theme: %v\n", err)
			os.Exit(1)
		}
	}

	var input io.Reader

	if flag.NArg() > 0 {
		file, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		input = file
	} else {
		input = os.Stdin
	}

	f := formatter.NewFormatter(themeManager)
	if err := f.ProcessMarkdown(input); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing Markdown: %v\n", err)
		os.Exit(1)
	}
}
