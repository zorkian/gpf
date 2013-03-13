/*******
 * gpf-server.go - a Go pf/gf server.
 * written by Mark Smith <mark@qq.is>
 */

package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	DIR string = "./files/"
)

type GpfState struct {
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	gpf := &GpfState{}

	srvr := &http.Server{
		Addr:         ":3466",
		Handler:      gpf,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(srvr.ListenAndServe())
}

// This is so sad. Fix it.
func randomString() string {
	return fmt.Sprintf("%d", rand.Int63())
}

func (self *GpfState) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "PUT" {
		if req.RequestURI != "/up" {
			http.NotFound(rw, req)
			return
		}
		fn := randomString()
		f, err := os.Create(filepath.Join(DIR, fn))
		if err != nil {
			http.Error(rw, "Failed to write file", 500)
			return
		}
		if _, err = io.Copy(f, req.Body); err != nil {
			http.Error(rw, "Failed to write file [2]", 500)
			return
		}
		if err = f.Close(); err != nil {
			http.Error(rw, "Failed to write file [3]", 500)
			return
		}
		fmt.Fprintf(rw, fn)
		return
	} else if req.Method != "GET" {
		http.NotFound(rw, req)
		return
	}

	f, err := os.Open(filepath.Join(DIR, filepath.Clean(req.RequestURI)))
	if err != nil {
		http.NotFound(rw, req)
		return
	}
	defer f.Close()

	http.ServeContent(rw, req, "filename", time.Now(), f)
}
