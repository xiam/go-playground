// Copyright 2011-2016 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
)

var (
	flagCompileURL   = flag.String("s", "https://play.golang.org/compile?output=json", "Sandbox service URL.")
	flagDisableCache = flag.Bool("disable-cache", false, "Disable cache.")
)

func init() {
	http.HandleFunc("/compile", compileHandler)
}

func compileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Origin", *flagAllowOriginHeader)

	if err := passThru(w, r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("compile error: %q", err)
		fmt.Fprintln(w, "Compile server error.")
	}
}

func passThru(w io.Writer, req *http.Request) error {
	var body bytes.Buffer

	_, err := io.Copy(&body, io.LimitReader(req.Body, maxSnippetSize+1))
	req.Body.Close()

	if err != nil {
		return fmt.Errorf("Error reading body: %q", err)
	}
	if body.Len() > maxSnippetSize {
		return fmt.Errorf("Snippet is too large")
	}

	snip := &Snippet{Body: body.Bytes()}
	id := snip.Id()
	key := []byte(id)

	var output bytes.Buffer

	if err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketCache)
		data := b.Get(key)
		if data == nil || *flagDisableCache {
			client := http.Client{}
			r, err := client.Post(*flagCompileURL, req.Header.Get("Content-Type"), &body)
			if err != nil {
				return err
			}
			defer r.Body.Close()

			data, err = ioutil.ReadAll(io.LimitReader(r.Body, maxSnippetSize+1))
			if len(data) > maxSnippetSize {
				return fmt.Errorf("Output is too large.")
			}
			if err = b.Put(key, data); err != nil {
				return err
			}
		}
		output.Write(data)
		return nil
	}); err != nil {
		return err
	}

	if _, err := io.Copy(w, &output); err != nil {
		return err
	}

	return nil
}
