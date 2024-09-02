package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/et0x/md2term/internal/formatter"
	"github.com/et0x/md2term/internal/theme"
)

// Version is the current version of md2term
const Version = "1.0.2"

func main() {
	// Parse command-line flags
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	themeFlag := flag.String("theme", "", "Specify the theme to use (default, dark, or light)")
	setDefaultThemeFlag := flag.Bool("set-default-theme", false, "Set the specified theme as the default")
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
		themeName, err = loadDefaultTheme()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading default theme: %v\n", err)
			themeName = "default" // Fallback to default theme
		}
	}

	err := themeManager.LoadTheme(themeName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading theme: %v\n", err)
		os.Exit(1)
	}

	// Set default theme if flag is set
	if *setDefaultThemeFlag {
		if err := saveDefaultTheme(themeName); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving default theme: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Default theme set to: %s\n", themeName)
		return
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

func loadDefaultTheme() (string, error) {
	configDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting user home directory: %w", err)
	}

	configFile := filepath.Join(configDir, ".md2term", "config")
	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "default", nil
		}
		return "", fmt.Errorf("error reading config file: %w", err)
	}

	return string(data), nil
}

func saveDefaultTheme(themeName string) error {
	configDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting user home directory: %w", err)
	}

	md2termDir := filepath.Join(configDir, ".md2term")
	err = os.MkdirAll(md2termDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	configFile := filepath.Join(md2termDir, "config")
	err = os.WriteFile(configFile, []byte(themeName), 0644)
	if err != nil {
		return fmt.Errorf("error saving config file: %w", err)
	}

	return nil
}
