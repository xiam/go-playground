// Copyright 2011-2016 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"go/format"
	"net/http"

	"golang.org/x/tools/imports"
)

func init() {
	http.HandleFunc("/fmt", fmtHandler)
}

type fmtResponse struct {
	Body  string
	Error string
}

func fmtHandler(w http.ResponseWriter, r *http.Request) {
	var (
		in  = []byte(r.FormValue("body"))
		out []byte
		err error
	)

	if len(in) > maxSnippetSize {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Compile server error.")
		return
	}

	if r.FormValue("imports") != "" {
		out, err = imports.Process("prog.go", in, nil)
	} else {
		out, err = format.Source(in)
	}

	var resp fmtResponse
	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.Body = string(out)
	}

	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Origin", *flagAllowOriginHeader)

	json.NewEncoder(w).Encode(resp)
}
