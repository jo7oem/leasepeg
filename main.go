package main

import (
	rest2 "exp-myddns/logic/rest"
	"exp-myddns/server/rest"
	"log/slog"
	"net/http"
	"sync"
)

func main() {
	logger := slog.Default().Handler()
	wg := sync.WaitGroup{}
	wg.Go(func() {
		http.ListenAndServe(":18888", rest.NewRegisterHandler(logger))
	})
	wg.Go(func() {
		http.ListenAndServe(":18889", rest2.NewHandler(logger))
	},
	)

	wg.Wait()

}
