package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jamesbee/srv/config"
)

var (
	cfgFile string
	rootCmd = cobra.Command{
		Use:   "srv",
		Short: "Serve your file system at once.",
		Long:  "James Bacon <jambecome@gmail.com>.\nServe your file or directory through http.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolP("help", "h", false, "Show this message")
	rootCmd.PersistentFlags().BoolVarP(&config.Debug, "debug", "d", false, "Enable debug output.")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.srv.yaml)")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")

	viper.BindPFlag("viper", rootCmd.PersistentFlags().Lookup("viper"))

	rootCmd.InitDefaultHelpCmd()
	rootCmd.InitDefaultVersionFlag()

	rootCmd.AddCommand(upCmd)
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".srv")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
