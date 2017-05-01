// Copyright 2011-2016 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
)

const (
	maxSnippetSize = 64 * 1024
)

var (
	flagDatabaseFile = flag.String("db", "playground.db", "Database file.")
	flagAllowShare   = flag.Bool("allow-share", false, "Allow storing and sharing snippets.")
)

var (
	db             *bolt.DB
	bucketSnippets = []byte("snippets")
	bucketCache    = []byte("cache")
	bucketConfig   = []byte("config")
	salt           []byte
)

type Snippet struct {
	Body []byte
}

func (s *Snippet) Id() string {
	h := sha1.New()
	h.Write(salt)
	h.Write(s.Body)
	sum := h.Sum(nil)
	return fmt.Sprintf("%x", sum)
}

func createBucket(name []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(name)
		return err
	})
}

func shareHandler(w http.ResponseWriter, r *http.Request) {
	if !*flagAllowShare || r.Method != "POST" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var body bytes.Buffer
	_, err := io.Copy(&body, io.LimitReader(r.Body, maxSnippetSize+1))
	r.Body.Close()
	if err != nil {
		log.Printf("Error reading body: %q", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	if body.Len() > maxSnippetSize {
		http.Error(w, "Snippet is too large", http.StatusRequestEntityTooLarge)
		return
	}

	snip := &Snippet{Body: body.Bytes()}
	id := snip.Id()
	key := []byte(id)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketSnippets)
		return b.Put(key, snip.Body)
	})

	if err != nil {
		log.Printf("Error putting snippet: %q", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Origin", *flagAllowOriginHeader)

	fmt.Fprint(w, id)
}
