package cmd

import (
	"github.com/spf13/cobra"
	"go_grpc_demo/server/impl"
)

// represents the serve command
func init() {
	var (
		cfgFile  = ""
		printCfg = false

		serveCmd = &cobra.Command{
			Use:   "serve",
			Short: "start server",
			Run: func(cmd *cobra.Command, args []string) {
				impl.StartServer(cfgFile, printCfg)
			},
		}
	)

	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&cfgFile, "config", "c", "server.toml", "config filename")
	serveCmd.Flags().BoolVarP(&printCfg, "print-config", "", false, "print config contents")
}
