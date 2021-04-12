package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"
)

func Run(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// First, find out if we even have a zig executable in our path.
	zigExe, err := exec.LookPath("zig")
	if err != nil {
		http.Error(w, "no zig compiler found", http.StatusInternalServerError)
		return
	}

	// Set up the temporary resources.
	dir, err := ioutil.TempDir("", "playground")
	if err != nil {
		http.Error(w, "creating temporary directory", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(dir)

	tmpSource := filepath.Join(dir, "play.zig")
	zigSource, err := ioutil.ReadAll(r.Body)
        if err != nil {
          http.Error(w, "reading body", http.StatusInternalServerError)
        }

	if err := ioutil.WriteFile(tmpSource, []byte(zigSource), 0666); err != nil {
		http.Error(w, "copying zig source", http.StatusInternalServerError)
		return
	}

	// Currently we cap compilation times at ten seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	output, err := exec.CommandContext(ctx, zigExe, "run", tmpSource).CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	_, err = w.Write(output)
	if err != nil {
		http.Error(w, "writing response", http.StatusInternalServerError)
	}
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

	router := httprouter.New()
	router.ServeFiles("/*filepath", http.Dir("static"))
	router.POST("/server/run", Run)

	chain := alice.New(rlMiddle.Handle, handlers.CompressHandler).Then(router)
	log.Fatal(http.ListenAndServe(":8080", chain))
}
