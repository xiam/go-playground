// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	bindAddr = ":3003"
)

var (
	allowedEnvs = []string{
		"GOOS",
		"GOARCH",
		"GOPATH",
		"GOCACHE",
		"TMPDIR",
		"GOTMPDIR",
	}
)

const maxRunTime = 5 * time.Second

type Request struct {
	Body string
}

type Response struct {
	Errors string
	Events []Event
}

func main() {
	http.HandleFunc("/compile", compileHandler)
	http.HandleFunc("/status", healthHandler)

	for _, envName := range allowedEnvs {
		if os.Getenv(envName) != "" {
			log.Printf("%s=%s", envName, os.Getenv(envName))
		}
	}

	log.Printf("Listening on %s...", bindAddr)

	var errCh = make(chan error)
	go func() {
		errCh <- http.ListenAndServe(bindAddr, nil)
	}()
	err := <-errCh

	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting server: %v", err)
		}
	}
}

func compileHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	version := r.PostFormValue("version")
	if version == "2" {
		req.Body = r.PostFormValue("body")
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("error decoding request: %v", err), http.StatusBadRequest)
			return
		}
	}
	resp, err := compileAndRun(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func compileAndRun(req *Request) (*Response, error) {
	tmpDir, err := os.MkdirTemp("", "sandbox")
	if err != nil {
		return nil, fmt.Errorf("error creating temp directory: %v", err)
	}

	defer os.RemoveAll(tmpDir)

	in := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(in, []byte(req.Body), 0400); err != nil {
		return nil, fmt.Errorf("error creating temp file %q: %v", in, err)
	}

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, in, nil, parser.PackageClauseOnly)
	if err == nil && f.Name.Name != "main" {
		return &Response{Errors: "package name must be main"}, nil
	}

	exe := filepath.Join(tmpDir, "a.out")
	log.Printf("Compiling %s to %s", in, exe)

	cmd := exec.Command("go", "build", "-o", exe, in)
	cmd.Dir = tmpDir

	cmd.Env = []string{}
	for _, envName := range allowedEnvs {
		cmd.Env = append(cmd.Env, envName+"="+os.Getenv(envName))
	}

	if out, err := cmd.CombinedOutput(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// Return compile errors to the user.

			// Rewrite compiler errors to refer to 'prog.go'
			// instead of '/tmp/sandbox1234/main.go'.
			errs := strings.Replace(string(out), in, "prog.go", -1)

			// "go build", invoked with a file name, puts this odd
			// message before any compile errors; strip it.
			errs = strings.Replace(errs, "# command-line-arguments\n", "", 1)

			return &Response{Errors: errs}, nil
		}
		return nil, fmt.Errorf("error building go source: %v", err)
	}

	cmd = exec.Command(exe)
	rec := new(Recorder)

	cmd.Stdout = rec.Stdout()
	cmd.Stderr = rec.Stderr()

	if err := runTimeout(cmd, maxRunTime); err != nil {
		if err == timeoutErr {
			return &Response{Errors: "process took too long"}, nil
		}
		if _, ok := err.(*exec.ExitError); !ok {
			return nil, fmt.Errorf("error running sandbox: %v", err)
		}
	}

	events, err := rec.Events()
	if err != nil {
		return nil, fmt.Errorf("error decoding events: %v", err)
	}

	return &Response{Events: events}, nil
}

var timeoutErr = errors.New("process timed out")

func runTimeout(cmd *exec.Cmd, d time.Duration) error {
	command := strings.Join(cmd.Args, " ")

	start := time.Now()
	log.Printf("Running %s", command)
	defer func() {
		elapsed := time.Since(start)
		log.Printf("Finished running %s, took %s", command, elapsed)
	}()

	if err := cmd.Start(); err != nil {
		return err
	}

	errc := make(chan error, 1)

	go func() {
		errc <- cmd.Wait()
	}()

	t := time.NewTimer(d)
	select {
	case err := <-errc:
		t.Stop()
		return err
	case <-t.C:
		cmd.Process.Kill()
		return timeoutErr
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if err := healthCheck(); err != nil {
		http.Error(w, "Health check failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "ok")
}

func healthCheck() error {
	resp, err := compileAndRun(&Request{Body: healthProg})
	if err != nil {
		return err
	}

	if resp.Errors != "" {
		return fmt.Errorf("compile error: %v", resp.Errors)
	}

	if len(resp.Events) != 1 || resp.Events[0].Message != "ok" {
		return fmt.Errorf("unexpected output: %v", resp.Events)
	}

	return nil
}

const healthProg = `
package main

import "fmt"

func main() { fmt.Print("ok") }
`
