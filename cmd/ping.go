package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "自分のIPアドレスを取得する",
	RunE:  runPing,
}

func runPing(cmd *cobra.Command, args []string) error {
	serverURL, _ := cmd.Flags().GetString("server-url")
	timeout, _ := cmd.Flags().GetDuration("timeout")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", serverURL+"/ping", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("リクエスト失敗: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		SrcAddr string `json:"src_addr"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("レスポンス解析失敗: %w", err)
	}

	fmt.Printf("あなたのIP: %s\n", result.SrcAddr)
	return nil
}
