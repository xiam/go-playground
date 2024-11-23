// Copyright 2011-2016 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/rand"
	"flag"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	flagDisableCache = flag.Bool("z", false, "Disable cache.")
	flagHelp         = flag.Bool("h", false, "Show help.")

	flagListenAddr        = flag.String("l", ":3000", "Listen address.")
	flagCompileURL        = flag.String("c", "https://play.golang.org/compile?output=json", "Compiler service URL.")
	flagAllowOriginHeader = flag.String("o", "*", "Access-Control-Allow-Origin header.")
)

func main() {
	flag.Parse()

	if *flagHelp {
		flag.PrintDefaults()
		return
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", createHandler)

	r.Handle("/*", http.FileServer(http.Dir("./static")))

	r.Post("/fmt", fmtHandler)
	r.Post("/compile", compileHandler)

	if *flagAllowShare {
		var err error

		if db, err = bolt.Open(*flagDatabaseFile, 0600, nil); err != nil {
			log.Fatal(err)
		}

		if err = createBucket(bucketSnippets); err != nil {
			log.Fatal(err)
		}

		if err = createBucket(bucketConfig); err != nil {
			log.Fatal(err)
		}

		if err = createBucket(bucketCache); err != nil {
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

		r.Get("/p/*", editHandler)
		r.Post("/share", shareHandler)
	}

	log.Printf("Serving Go Playground at %v\n", *flagListenAddr)
	http.ListenAndServe(*flagListenAddr, r)
}
