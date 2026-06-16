package rest

import (
	"context"
	"fmt"
	"net"
	"net/http"
)
import (
	oapi "exp-myddns/gateway/rest"
)

func NewHandler() http.Handler {
	return oapi.Handler(oapi.NewStrictHandler(new(handler), []oapi.StrictMiddlewareFunc{strictIPRecordHandler}))
}

type handler struct {
}

func (h handler) GetHealthcheck(ctx context.Context, request oapi.GetHealthcheckRequestObject) (oapi.GetHealthcheckResponseObject, error) {
	return oapi.GetHealthcheck200JSONResponse{Status: "OK"}, nil
}

func (h handler) GetPing(ctx context.Context, request oapi.GetPingRequestObject) (oapi.GetPingResponseObject, error) {
	srcIP, err := srcIPFromContext(ctx)
	if err != nil {
		srcIP = "pong"
	}

	resp := oapi.GetPing200JSONResponse{
		Ping: srcIP,
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
