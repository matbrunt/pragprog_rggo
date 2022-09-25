/*
Copyright © 2022 Mat Brunt
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags
	cfgFile string

	rootCmd = &cobra.Command{
		Version: "0.1",
		Use:     "pScan",
		Short:   "Fast TCP port scanner",
		Long: `pScan - short for Port Scanner - executes TCP port scan on a list of hosts.
		
		pScan allows you to add, list, and delete hosts from the list.
		
		pScan executes a port scan on specified TCP ports. You can customise the target ports using a command line flag.`,
		// Uncomment the following line if your bare application
		// has an uuusociated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}
)

// rootCmd represents the base command when called without any subcommands

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// https://github.com/spf13/cobra/blob/main/user_guide.md
	// Here you will define your flags and configuration settings.
	cobra.OnInitialize(initConfig)

	// Viper reads in environment variables starting with PSCAN_
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("PSCAN")

	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// persistent flags are available to the command, and all subcommands
	// under that command. By adding a persistent flag to the root command,
	// you make it global.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pScan.yaml)")

	// '--hosts-file' or '-f' flag, not specified defaults to 'pScan.hosts'
	rootCmd.PersistentFlags().StringP("hosts-file", "f", "pScan.hosts", "pScan hosts file")

	// allow user to override hosts-file config variable with environment variable
	// using PSCAN_HOSTS_FILE
	viper.BindPFlag("hosts-file", rootCmd.PersistentFlags().Lookup("hosts-file"))

	// rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	// viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))

	versionTemplate := `{{printf "%s: %s - version %s\n" .Name .Short .Version}}`
	rootCmd.SetVersionTemplate(versionTemplate)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".pScan" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".pScan")
	}

	// read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
