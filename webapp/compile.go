package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

var flagCompileURL = flag.String("c", "https://play.golang.org/compile?output=json", "Services prefix.")

func init() {
	http.HandleFunc("/compile", compileHandler)
}

func compileHandler(w http.ResponseWriter, r *http.Request) {
	if err := passThru(w, r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("compile error: %q", err)
		fmt.Fprintln(w, "Compile server error.")
	}
}

func passThru(w io.Writer, req *http.Request) error {
	client := http.Client{}
	defer req.Body.Close()
	r, err := client.Post(*flagCompileURL, req.Header.Get("Content-Type"), req.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if _, err := io.Copy(w, r.Body); err != nil {
		return err
	}
	return nil
}
