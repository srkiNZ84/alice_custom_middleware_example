package main

import (
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store"
)

func timeoutHandler(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, 1*time.Second, "timed out")
}

func myStripPrefix(h http.Handler) http.Handler {
	return http.StripPrefix("/old", h)
}

func myRandomResponse(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}

func myApp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func main() {
	th := throttled.RateLimit(throttled.PerSec(10), &throttled.VaryBy{Path: true}, store.NewMemStore(1000))
	myHandler := http.HandlerFunc(myApp)

	chain := alice.New(th.Throttle, myRandomResponse, myStripPrefix, timeoutHandler, nosurf.NewPure).Then(myHandler)
	http.ListenAndServe(":8000", chain)
}
