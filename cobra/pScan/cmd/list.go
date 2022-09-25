/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"pragprog.com/rggo/cobra/pScan/scan"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List hosts in hosts list",
	// Allow users to specify 'list' or 'l' as subcommand
	Aliases: []string{"l"},
	// Run executes the function when command invoked, but doesn't return an error
	// RunE returns error that's displayed to user if needed
	RunE: func(cmd *cobra.Command, args []string) error {
		// we only parse command args in here then call out to external function
		// to make testing easier
		hostsFile := viper.GetString("hosts-file")

		return listAction(os.Stdout, hostsFile, args)
	},
}

func listAction(out io.Writer, hostsFile string, args []string) error {
	// create an instance of the list
	hl := &scan.HostsList{}

	// loads content of file into list
	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	// iterate over each line, printing to writer as a new line
	for _, h := range hl.Hosts {
		if _, err := fmt.Fprintln(out, h); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	hostsCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
