package rest_test

import (
	gwRest "exp-myddns/gateway/rest"
	"exp-myddns/logic/rest"
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_handler_GetHealthcheckLive(t *testing.T) {
	tests := []struct {
		name       string
		wantStatus string
		wantErr    bool
	}{
		{
			name:       "success",
			wantStatus: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := rest.NewHandler(slog.DiscardHandler)
			srv := httptest.NewServer(h)
			defer srv.Close()
			c := srv.Client()
			client, err := gwRest.NewClient(srv.URL, gwRest.WithHTTPClient(c))
			if err != nil {
				t.Fatal(err)
			}

			resp, err := client.GetHealthcheckLive(t.Context())
			if tt.wantErr != (err != nil) {
				t.Errorf("GetHealthcheckLive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			hcResp, err := gwRest.ParseGetHealthcheckLiveResponse(resp)
			if err != nil {
				t.Errorf("GetHealthcheckLive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.wantStatus, hcResp.JSON200.Status); diff != "" {
				t.Errorf("GetHealthcheckLive() diff = %v", diff)
			}
		})
	}
}
