package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	headerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true).Underline(true)
)

func printSuccess(format string, args ...any) {
	fmt.Println(successStyle.Render(fmt.Sprintf(format, args...)))
}

func printError(format string, args ...any) {
	fmt.Println(errorStyle.Render(fmt.Sprintf(format, args...)))
}

func printInfo(format string, args ...any) {
	fmt.Println(infoStyle.Render(fmt.Sprintf(format, args...)))
}

func printHeader(format string, args ...any) {
	fmt.Println(headerStyle.Render(fmt.Sprintf(format, args...)))
}
