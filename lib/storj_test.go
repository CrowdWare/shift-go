package lib

import (
	"bytes"
	"context"
	"log"
	"testing"
)

func TestStorj(t *testing.T) {
	Init("/tmp")
	err := initStorj(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	uploadBuffer := []byte("one fish two fish red fish blue fish")
	err = put("foo/bar/baz", uploadBuffer)
	if err != nil {
		log.Fatal(err)
	}

	buffer, err := get("foo/bar/baz")
	if err != nil {
		log.Fatal(err)
	}

	if !bytes.Equal(uploadBuffer, buffer) {
		t.Error("Storj buffers are not identical")
	}

	err = delete("foo/bar/baz")
	if err != nil {
		log.Fatal(err)
	}
}
