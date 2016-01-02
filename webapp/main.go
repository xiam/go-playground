// Copyright 2011-2016 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	flagListenAddr = flag.String("l", ":3000", "Listen address.")
	flagHelp       = flag.Bool("h", false, "Show help.")
)

func main() {
	flag.Parse()

	if *flagHelp {
		flag.PrintDefaults()
		return
	}

	log.Printf("Serving Go playground at %v...\n", *flagListenAddr)

	if err := http.ListenAndServe(*flagListenAddr, nil); err != nil {
		log.Fatal(err)
	}
}
