package rest

import (
	"context"
	"errors"
	"fmt"
	"leasepeg/logging"
	"leasepeg/store/memory"
	"log/slog"
	"net"
	"net/http"
)
import (
	oapi "leasepeg/internal/generated/oapi/rest"
)

func NewHandler(lh slog.Handler, opts ...logging.LoggerOption) http.Handler {
	logger := logging.NewLogger(lh, "restHandle", opts...)
	return oapi.Handler(oapi.NewStrictHandler(new(handler{
		logger: logger,
	}), []oapi.StrictMiddlewareFunc{strictIPRecordHandler}))
}

type handler struct {
	logger *logging.Logger
}

func (h handler) GetHealthcheckLive(_ context.Context, _ oapi.GetHealthcheckLiveRequestObject) (oapi.GetHealthcheckLiveResponseObject, error) {
	return oapi.GetHealthcheckLive200JSONResponse{
		Status: "ok",
	}, nil
}

func (h handler) GetHealthcheckReady(ctx context.Context, request oapi.GetHealthcheckReadyRequestObject) (oapi.GetHealthcheckReadyResponseObject, error) {
	state := memory.GetHealthCheckState()
	var msg *string
	if state.Err != nil && state.Err.Error() != "" {
		msg = new(state.Err.Error())
	}

	res := oapi.HealthcheckReadyResp{
		Message: msg,
		Status:  oapi.HealthcheckStatus(state.State.String()),
	}

	switch state.State {
	case memory.StateReady, memory.StateDegraded:
		return oapi.GetHealthcheckReady200JSONResponse(res), nil
	case memory.StateStarting, memory.StateShuttingDown:
		return oapi.GetHealthcheckReady503JSONResponse(res), nil
	case memory.StateError:
		return oapi.GetHealthcheckReady500JSONResponse(res), nil
	default:
		return nil, errors.New("unknown health check state : " + state.State.String())
	}

}

func (h handler) GetPing(ctx context.Context, request oapi.GetPingRequestObject) (oapi.GetPingResponseObject, error) {
	srcAddr, err := srcIPFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get src IP : %w", err)
	}

	srcIP := srcAddr
	if host, _, err := net.SplitHostPort(srcAddr); err == nil {
		srcIP = host
	}

	resp := oapi.GetPing200JSONResponse{
		SrcAddr: srcIP,
	}

	return resp, nil
}

type srcIPKeyType struct{}

func srcIPFromContext(ctx context.Context) (string, error) {
	srcIP, ok := ctx.Value(srcIPKeyType{}).(string)
	if !ok {
		return "", fmt.Errorf("src IP not found in context")
	}

	return srcIP, nil
}

func withSrcIP(ctx context.Context, srcIP string) context.Context {
	return context.WithValue(ctx, srcIPKeyType{}, srcIP)
}

func strictIPRecordHandler(f oapi.StrictHandlerFunc, operationID string) oapi.StrictHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}

		ctx = withSrcIP(ctx, host)
		return f(ctx, w, r, request)
	}
}
