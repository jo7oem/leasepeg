package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health [live|ready]",
	Short: "サーバーのヘルスチェックを行う",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runHealth,
}

func init() {
	healthCmd.Flags().String("type", "ready", "チェックタイプ (live|ready)")
}

func runHealth(cmd *cobra.Command, args []string) error {
	serverURL, _ := cmd.Flags().GetString("server-url")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	checkType, _ := cmd.Flags().GetString("type")

	if len(args) > 0 {
		checkType = args[0]
	}

	path := "/healthcheck/ready"
	if checkType == "live" {
		path = "/healthcheck/live"
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", serverURL+path, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("リクエスト失敗: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("ステータス: %d\n", resp.StatusCode)
	return nil
}
