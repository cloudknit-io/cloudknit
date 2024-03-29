package cmd

import (
	"fmt"
	"os"

	"github.com/compuzest/zlifecycle-internal-cli/app/cmd/aws"
	"github.com/compuzest/zlifecycle-internal-cli/app/cmd/git"
	"github.com/compuzest/zlifecycle-internal-cli/app/cmd/state"

	"github.com/compuzest/zlifecycle-internal-cli/app/common"
	"github.com/compuzest/zlifecycle-internal-cli/app/env"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "zli [command]",
		Version: env.Version,
		Short:   "zLifecycle internal CLI",
		Long:    `zLifecycle internal CLI for administrative management and workflow executor`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				if err := cmd.Help(); err != nil {
					common.Failure(1)
				}
				common.Success()
			}
		},
	}

	cobra.OnInitialize(initConfig)

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.zl.yaml)")
	cmd.PersistentFlags().BoolVarP(&env.Verbose, "verbose", "v", false, "enable command logs")
	cmd.AddCommand(git.NewRootCmd())
	cmd.AddCommand(state.NewRootCmd())
	cmd.AddCommand(aws.NewRootCmd())

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(NewRootCmd().Execute())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".zl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".zl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
