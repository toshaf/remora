package server

import (
	"os"
	"path"
	"testing"
)

func Test_Paths_Correct(t *testing.T) {
	wd, err := os.Getwd()
	Check(err)

	for _, c := range []struct {
		in     string
		target string
		pipes  string
	}{
		// relative path
		{
			in:     "path/to/a/file",
			target: path.Join(wd, "path/to/a/file"),
			pipes:  path.Join(wd, "path/to/a/.pipes"),
		},
		// absolute path
		{
			in:     "/path/to/a/file",
			target: "/path/to/a/file",
			pipes:  "/path/to/a/.pipes",
		},
	} {
		srv := New(Args{
			Target: c.in,
		}).(*server)

		if srv.target != c.target {
			t.Errorf("Expected target %s but got %s", c.target, srv.target)
		}

		if srv.pipes != c.pipes {
			t.Errorf("Expected pipes %s but got %s", c.pipes, srv.pipes)
		}
	}
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}
