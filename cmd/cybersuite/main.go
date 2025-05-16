package main

import (
	"fmt"
	"os"

	"suite/internal/tools/portscan"
	"suite/pkg/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	// Initialize root command
	rootCmd := &cobra.Command{
		Use:   "cybersuite",
		Short: "A comprehensive cybersecurity tool suite",
		Long: `A terminal-based suite of cybersecurity tools including:
- Port Scanner
- Subdomain Enumerator
- Hash Identifier & Cracker
- Network Mapper
- HTTP Header Security Scanner`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := startTUI(); err != nil {
				fmt.Printf("Error starting TUI: %v\n", err)
				os.Exit(1)
			}
		},
	}

	// Add flags for direct tool execution
	rootCmd.PersistentFlags().StringP("tool", "t", "", "Tool to run (portscan, subdomain, hashcrack, netmap, headers)")
	rootCmd.PersistentFlags().StringP("target", "T", "", "Target to scan")

	// Initialize config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.cybersuite")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Error reading config: %v\n", err)
			os.Exit(1)
		}
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
}

func startTUI() error {
	p := tea.NewProgram(ui.New())
	model, err := p.Run()
	if err != nil {
		return err
	}

	// Handle tool selection
	mainModel := model.(ui.Model)
	switch mainModel.CurrentTool() {
	case "Port Scanner":
		p := tea.NewProgram(portscan.NewModel())
		if _, err := p.Run(); err != nil {
			return err
		}
	}

	return nil
}
