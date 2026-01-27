package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

type Command int
type CtxKey string

const (
	R Command = iota
	F
)

const CtxLogger CtxKey = "logger"

func LoggingMiddleware(h http.Handler, logger zerolog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), CtxLogger, logger)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func whichZig(r *http.Request) string {
	const zigVersion = "X-Zig-Version"

	return r.Header.Get(zigVersion)
}

func execute(w http.ResponseWriter, r *http.Request, command Command) {
	logger := r.Context().Value(CtxLogger).(zerolog.Logger)

	// Limit how big a source file can be. 5MB here.
	r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)
	zigSource, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error().Err(err).Msg("reading the request body")
		http.Error(w, "reading body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Set up the temporary resources.
	playgroundDir := os.Getenv("PLAYGROUND_DIR")
	dir, err := os.MkdirTemp(playgroundDir, "playground")
	if err != nil {
		logger.Error().Err(err).Msg("making the temporary directory")
		http.Error(w, "creating temporary directory", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(dir)

	tmpSource := filepath.Join(dir, "play.zig")
	if err := os.WriteFile(tmpSource, []byte(zigSource), 0666); err != nil {
		logger.Error().Err(err).Msg("copying the source")
		http.Error(w, "copying zig source", http.StatusInternalServerError)
		return
	}

	// Currently we cap compilation times at thirty seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// We only have two commands for now.
	var output []byte
	if command == R {
		home := os.Getenv("HOME")
		output, err = exec.CommandContext(ctx, path.Join(home, "zrun.sh"), whichZig(r), tmpSource).CombinedOutput()
	} else {
		fd, ferr := os.Open(tmpSource)
		if ferr != nil {
			logger.Error().Err(err).Msg("opening the tmp source during fmt command")
		}
		cmd := exec.CommandContext(ctx, "zvm", "run", whichZig(r), "fmt", "--stdin")
		cmd.Stdin = fd
		output, err = cmd.CombinedOutput()
		output = output[:len([]byte(zigSource))]
	}

	if err != nil {
		logger.Error().Err(err).Msg("running the command")
		w.WriteHeader(http.StatusBadRequest)
	}

	_, err = w.Write(output)
	if err != nil {
		logger.Error().Err(err).Msg("writing the response")
		http.Error(w, "writing response", http.StatusInternalServerError)
	}
}

func Run(w http.ResponseWriter, r *http.Request) {
	execute(w, r, R)
}

func Fmt(w http.ResponseWriter, r *http.Request) {
	execute(w, r, F)
}
