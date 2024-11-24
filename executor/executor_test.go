package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompileAndRun(t *testing.T) {
	var tests = []struct {
		prog, want, errors string
	}{
		{prog: `
package main

import "time"

func main() {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err.Error())
	}
	println(loc.String())
}
`, want: "America/New_York"},
		{prog: `
package test

func main() {
    println("test")
}
`, want: "", errors: "package name must be main"},
	}

	for _, tt := range tests {
		resp, err := compileAndRun(&Request{Body: tt.prog})
		require.NoError(t, err)

		if tt.errors != "" {
			assert.Equal(t, tt.errors, resp.Errors)
			continue
		}

		assert.Equal(t, "", resp.Errors)
		assert.Equal(t, 1, len(resp.Events))

		assert.Contains(t, resp.Events[0].Message, tt.want)
	}
}
