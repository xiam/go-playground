// Copyright 2011-2016 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/boltdb/bolt"
)

func init() {
	http.HandleFunc("/", editHandler)
}

var editTemplate = template.Must(template.ParseFiles("static/index.html"))

type editData struct {
	Snippet *Snippet
	Share   bool
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	snip := &Snippet{Body: []byte(helloPlayground)}
	if strings.HasPrefix(r.URL.Path, "/p/") {
		if !*flagAllowShare {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		id := r.URL.Path[3:]
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(bucketSnippets)
			v := b.Get([]byte(id))
			if v != nil {
				snip.Body = v
			}
			return nil
		})
	}

	editTemplate.Execute(w, &editData{snip, *flagAllowShare})
}

const helloPlayground = `package main

import "fmt"

func main() {
	fmt.Println("Hello, playground")
}
`
