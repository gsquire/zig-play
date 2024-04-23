package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
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

func execute(w http.ResponseWriter, r *http.Request, command Command) {
	const defaultZig = "/usr/local/bin/zig"

	var zigExe string
	foundZigExe, zigExeErr := exec.LookPath("zig")
	if zigExeErr != nil {
		zigExe = defaultZig
	} else {
		zigExe = foundZigExe
	}

	logger := r.Context().Value(CtxLogger).(zerolog.Logger)

	// Limit how big a source file can be. 5MB here.
	r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)
	zigSource, err := ioutil.ReadAll(r.Body)
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
	if err := ioutil.WriteFile(tmpSource, []byte(zigSource), 0666); err != nil {
		logger.Error().Err(err).Msg("copying the source")
		http.Error(w, "copying zig source", http.StatusInternalServerError)
		return
	}

	// Currently we cap compilation times at ten seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// We only have two commands for now.
	var output []byte
	if command == R {
		output, err = exec.CommandContext(ctx, zigExe, "run", "--global-cache-dir", dir, tmpSource).CombinedOutput()
	} else {
		// The global cache directory option is not available for this command.
		cmd := fmt.Sprintf("ZIG_GLOBAL_CACHE_DIR=%s cat %s | %s fmt --stdin", dir, tmpSource, zigExe)
		output, err = exec.CommandContext(ctx, "bash", "-c", cmd).CombinedOutput()
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
	logger := r.Context().Value(CtxLogger).(zerolog.Logger)
	logger.Info().Msg("running compile")

	execute(w, r, R)
}

func Fmt(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value(CtxLogger).(zerolog.Logger)
	logger.Info().Msg("running format")

	execute(w, r, F)
}
