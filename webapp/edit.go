// Copyright 2011-2016 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"

	"github.com/boltdb/bolt"
)

func editHandler(w http.ResponseWriter, r *http.Request) {
	snip := &Snippet{Body: []byte(helloPlayground)}

	var id string

	if len(r.URL.Path) >= 3 {
		id = r.URL.Path[3:]
	}

	db.View(func(tx *bolt.Tx) error {
		if id == "" {
			return nil
		}
		b := tx.Bucket(bucketSnippets)
		v := b.Get([]byte(id))
		if v != nil {
			snip.Body = v
		}
		return nil
	})

	editorTemplate.Execute(w, &editorText{snip, *flagAllowShare})
}
