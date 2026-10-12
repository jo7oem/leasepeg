package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var clientCmd = &cobra.Command{
	Use:               "client",
	Short:             "DDNSクライアント操作",
	PersistentPreRunE: initClientConfig,
}

func init() {
	clientCmd.PersistentFlags().String("server-url", "http://localhost:8080", "サーバーURL (環境変数: LEASEPEG_CLIENT_SERVER_URL)")
	clientCmd.PersistentFlags().Duration("timeout", 10*time.Second, "リクエストタイムアウト (環境変数: LEASEPEG_CLIENT_TIMEOUT)")

	viper.BindPFlag("client.server_url", clientCmd.PersistentFlags().Lookup("server-url"))
	viper.BindPFlag("client.timeout", clientCmd.PersistentFlags().Lookup("timeout"))

	clientCmd.AddCommand(healthCmd)
	clientCmd.AddCommand(pingCmd)
}

func initClientConfig(cmd *cobra.Command, args []string) error {
	return nil
}
