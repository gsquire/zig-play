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
)

const (
	defaultZig  = "/usr/local/bin/zig"
	maxFileSize = 5 * 1024 * 1024
	maxCompTime = 10 * time.Second
)

func initZig() string {
	foundZigExe, zigExeErr := exec.LookPath("zig")
	if zigExeErr != nil {
		return defaultZig
	}
	return foundZigExe
}

func setupPlayground(w http.ResponseWriter, r *http.Request) (string, string, error) {
	// Set up the temporary resources.
	dir, err := ioutil.TempDir("", "playground")
	if err != nil {
		return "", "", err
	}

	// Limit how big a source file can be. 5MB here.
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)
	zigSource, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", "", err
	}

	tmpSource := filepath.Join(dir, "play.zig")
	if err := ioutil.WriteFile(tmpSource, []byte(zigSource), 0666); err != nil {
		return "", "", err
	}

	return dir, tmpSource, nil
}

func Run(w http.ResponseWriter, r *http.Request) {
	zigExe := initZig()
	defer r.Body.Close()

	dir, tmpSource, err := setupPlayground(w, r)
	if err != nil {
		http.Error(w, "reading source", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(dir)

	// Currently we cap compilation times at ten seconds.
	ctx, cancel := context.WithTimeout(context.Background(), maxCompTime)
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

func Fmt(w http.ResponseWriter, r *http.Request) {
	zigExe := initZig()
	defer r.Body.Close()

	dir, tmpSource, err := setupPlayground(w, r)
	if err != nil {
		http.Error(w, "reading source", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(dir)

	// Currently we cap format times at ten seconds.
	ctx, cancel := context.WithTimeout(context.Background(), maxCompTime)
	defer cancel()

	cmd := fmt.Sprintf("cat %s | %s fmt --stdin", tmpSource, zigExe)
	output, err := exec.CommandContext(ctx, "bash", "-c", cmd).CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	_, err = w.Write(output)
	if err != nil {
		http.Error(w, "writing response", http.StatusInternalServerError)
	}
}
