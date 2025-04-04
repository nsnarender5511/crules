package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"cursor++/internal/utils"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// DisplayFileTable displays a list of files in a formatted table
func DisplayFileTable(files []utils.FileInfo) {
	if len(files) == 0 {
		Plain("No files found.")
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Type", "Name", "Size", "Modified"})

	for i, file := range files {
		fileType := "file"
		if file.IsDir {
			fileType = "dir"
		}

		t.AppendRow(table.Row{
			i + 1,
			fileType,
			file.RelativePath,
			utils.FormatFileSize(file.Size),
			file.ModTime.Format("Jan 02, 2006 15:04:05"),
		})
	}

	t.SetStyle(table.StyleLight)
	t.Style().Color.Header = text.Colors{text.FgHiBlue}
	t.Render()

	Plain("")
	Plain("Total: %d items", len(files))
}

// PromptYesNo asks the user a yes/no question and returns the answer
func PromptYesNo(question string) bool {
	for {
		Prompt("%s (y/n): ", question)

		var response string
		fmt.Scanln(&response)

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}

		Warning("Please enter 'y' or 'n'")
	}
}

// PromptOptions presents a list of options and returns the selected index
func PromptOptions(prompt string, options []string) int {
	Plain(prompt)

	for i, option := range options {
		Plain("  %d. %s", i+1, option)
	}

	for {
		Prompt("Enter your choice (1-%d): ", len(options))

		var choice string
		fmt.Scanln(&choice)

		if num, err := strconv.Atoi(choice); err == nil && num >= 1 && num <= len(options) {
			return num - 1
		}

		Warning("Invalid choice. Please enter a number between 1 and %d", len(options))
	}
}

// PromptInput asks for a string input with optional validation
func PromptInput(prompt string, validator func(string) bool) string {
	return PromptInputWithDefault(prompt, "", validator)
}

// PromptInputWithDefault asks for a string input with a default value and optional validation
func PromptInputWithDefault(prompt string, defaultValue string, validator func(string) bool) string {
	reader := bufio.NewReader(os.Stdin)

	promptText := prompt
	if defaultValue != "" {
		promptText = fmt.Sprintf("%s [%s]", prompt, defaultValue)
	}

	for {
		Prompt("%s ", promptText)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Use default value if input is empty
		if input == "" && defaultValue != "" {
			input = defaultValue
		}

		if validator == nil || validator(input) {
			return input
		}

		Warning("Invalid input. Please try again.")
	}
}

// ValidateURL validates if the input looks like a URL
func ValidateURL(input string) bool {
	// Check for HTTP/HTTPS/Git protocols
	if strings.HasPrefix(input, "http://") ||
		strings.HasPrefix(input, "https://") ||
		strings.HasPrefix(input, "git://") {
		return true
	}

	// Check for SSH git format (git@host:path/repo.git)
	if strings.HasPrefix(input, "git@") && strings.Contains(input, ":") {
		return true
	}

	return false
}
