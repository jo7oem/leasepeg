package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"leasepeg/logging"
	"leasepeg/server/rest"
	"leasepeg/store/config"
	"leasepeg/store/memory"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "DDNSサーバーを起動する",
	RunE:  runServer,
}

func init() {
	serverCmd.Flags().String("host", "", "サーバーホスト (環境変数: LEASEPEG_SERVER_HOST)")
	serverCmd.Flags().Int("port", 0, "サーバーポート (環境変数: LEASEPEG_SERVER_PORT)")
	viper.BindPFlag("server.host", serverCmd.Flags().Lookup("host"))
	viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port"))
}

func runServer(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("設定読み込み失敗: %w", err)
	}

	var handler slog.Handler
	if cfg.Logging.Format == "text" {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: parseLogLevel(cfg.Logging.Level),
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: parseLogLevel(cfg.Logging.Level),
		})
	}
	logger := logging.NewLogger(handler, "server")

	memory.SetHealthCheckState(memory.StateStarting, nil)
	defer memory.SetHealthCheckState(memory.StateShuttingDown, nil)

	httpHandler := rest.NewRegisterHandler(logger.Handler())

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      httpHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		logger.Info("サーバー起動", slog.String("addr", srv.Addr))
		memory.SetHealthCheckState(memory.StateReady, nil)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		logger.Info("シャットダウン開始")
	case err := <-errCh:
		if err != nil {
			logger.Error("サーバーエラー", err)
			memory.SetHealthCheckState(memory.StateError, err)
			return err
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("強制終了", err)
		return err
	}

	logger.Info("正常終了")
	return nil
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
