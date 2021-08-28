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

func Run(w http.ResponseWriter, r *http.Request) {
	const defaultZig = "/usr/local/bin/zig"

	var zigExe string
	foundZigExe, zigExeErr := exec.LookPath("zig")
	if zigExeErr != nil {
		zigExe = defaultZig
	} else {
		zigExe = foundZigExe
	}

	// Limit how big a source file can be. 5MB here.
	r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)
	zigSource, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "reading body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Set up the temporary resources.
	dir, err := ioutil.TempDir("", "playground")
	if err != nil {
		http.Error(w, "creating temporary directory", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(dir)

	tmpSource := filepath.Join(dir, "play.zig")
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

func Fmt(w http.ResponseWriter, r *http.Request) {
	const defaultZig = "/usr/local/bin/zig"

	var zigExe string
	foundZigExe, zigExeErr := exec.LookPath("zig")
	if zigExeErr != nil {
		zigExe = defaultZig
	} else {
		zigExe = foundZigExe
	}

	// Limit how big a source file can be. 5MB here.
	r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)
	zigSource, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "reading body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Set up the temporary resources.
	dir, err := ioutil.TempDir("", "playground")
	if err != nil {
		http.Error(w, "creating temporary directory", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(dir)

	tmpSource := filepath.Join(dir, "play.zig")
	if err := ioutil.WriteFile(tmpSource, []byte(zigSource), 0666); err != nil {
		http.Error(w, "copying zig source", http.StatusInternalServerError)
		return
	}

	// Currently we cap format times at ten seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
