// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"crypto/rand"
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
	flagDatabaseFile = flag.String("db", "snippets.db", "Snippets database.")
	flagAllowShare   = flag.Bool("allow-share", false, "Allow storing and sharing snippets.")
)

var (
	db             *bolt.DB
	bucketSnippets = []byte("snippets")
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

func init() {
	var err error
	http.HandleFunc("/share", shareHandler)
	if db, err = bolt.Open(*flagDatabaseFile, 0600, nil); err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketSnippets)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketConfig)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket(bucketConfig)
		salt = b.Get([]byte("salt"))
		if salt == nil {
			salt = make([]byte, 30)
			if _, err = rand.Read(salt); err != nil {
				return err
			}
			b.Put([]byte("salt"), salt)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
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

	fmt.Fprint(w, id)
}
