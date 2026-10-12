package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version string
	rootCmd = &cobra.Command{
		Use:               "leasepeg",
		Short:             "ダイナミックDNS サーバー & クライアント",
		Long:              `DDNSサーバー起動とクライアント操作を行うCLIツール`,
		PersistentPreRunE: initConfig,
		Run: func(cmd *cobra.Command, args []string) {
			if err := serverCmd.RunE(serverCmd, args); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
	}
)

func Execute(v string) error {
	version = v
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "", "設定ファイルパス")
	rootCmd.PersistentFlags().String("log-level", "info", "ログレベル (debug|info|warn|error)")
	rootCmd.PersistentFlags().String("log-format", "json", "ログ形式 (json|text)")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(clientCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "バージョンを表示する",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("leasepeg %s\n", version)
	},
}

func initConfig(cmd *cobra.Command, args []string) error {
	viper.SetEnvPrefix("LEASEPEG")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if configFile, _ := cmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(".")
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			viper.AddConfigPath(filepath.Join(xdg, "leasepeg"))
		} else if home, _ := os.UserHomeDir(); home != "" {
			viper.AddConfigPath(filepath.Join(home, ".config", "leasepeg"))
		}
		viper.AddConfigPath("/etc/leasepeg")
	}

	_ = viper.ReadInConfig()
	return nil
}
