package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"
	"github.com/unrolled/secure"
)

func securitySettings() *secure.Secure {
	return secure.New(secure.Options{
		BrowserXssFilter:   true,
		ContentTypeNosniff: true,
		FrameDeny:          true,
		STSPreload:         true,
		STSSeconds:         31536000,
	})
}

func main() {
	// Users can compile code 5 times per minute.
	rateLimiter, err := memorystore.New(&memorystore.Config{
		Tokens:   5,
		Interval: time.Minute,
	})
	if err != nil {
		log.Fatal("error making rate limiter", err)
	}
	rlMiddle, err := httplimit.NewMiddleware(rateLimiter, httplimit.IPKeyFunc())
	if err != nil {
		log.Fatal(err)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	router := httprouter.New()

	router.Handler(http.MethodPost, "/server/run", http.HandlerFunc(Run))
	router.Handler(http.MethodPost, "/server/fmt", http.HandlerFunc(Fmt))

	chain := alice.New(rlMiddle.Handle, securitySettings().Handler, handlers.CompressHandler, handlers.RecoveryHandler()).Then(LoggingMiddleware(router, logger))
	log.Fatal(http.ListenAndServe(":8080", chain))
}
