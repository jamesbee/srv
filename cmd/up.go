package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jamesbee/srv/server"
)

var (
	upCmd = &cobra.Command{
		Use:   "up",
		Short: "Bring up http(s) server.",
		Long:  "Bring up http(s) server.",
		Run:   server.CommandUp,
	}
)

func init() {
	upCmd.Flags().BoolP("help", "h", false, "Show this message")
	upCmd.PersistentFlags().BoolVarP(&server.UseTLS, "tls-enable", "t", false, "Enable TLS")
	upCmd.PersistentFlags().StringVarP(&server.Cert, "tls-cert", "c", "", "TLS Cert path, e.g. ./cert.pem")
	upCmd.PersistentFlags().StringVarP(&server.Key, "tls-key", "k", "", "TLS Key path, e.g. ./cert.key")
	upCmd.PersistentFlags().StringVar(&server.Host, "host", "127.0.0.1", "Server listen host")
	upCmd.PersistentFlags().IntVarP(&server.Port, "port", "p", 3000, "Server listen port")
	upCmd.PersistentFlags().BoolVarP(&server.EnableMarkdown, "markdown", "m", true, "Enable markdown parse")

	viper.BindPFlag("tlsCert", upCmd.PersistentFlags().Lookup("tls-cert"))
	viper.BindPFlag("tlsKey", upCmd.PersistentFlags().Lookup("tls-key"))
	viper.BindPFlag("host", upCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", upCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("markdown", upCmd.PersistentFlags().Lookup("markdown"))

	viper.SetDefault("markdown", "true")
	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("port", "3000")
}
