package main

import (
	rest2 "exp-myddns/logic/rest"
	"exp-myddns/server/rest"
	"net/http"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Go(func() {
		http.ListenAndServe(":18888", rest.NewRegisterHandler())
	})
	wg.Go(func() {
		http.ListenAndServe(":18889", rest2.NewHandler())
	},
	)

	wg.Wait()

}
