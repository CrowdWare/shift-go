package lib

import (
	"bytes"
	"context"
	"testing"
)

func TestStorj(t *testing.T) {
	Init("/tmp")
	err := initStorj(context.Background())
	if err != nil {
		t.Errorf("initStorj failed %s", err.Error())
		return
	}
	uploadBuffer := []byte("one fish two fish red fish blue fish")
	err = put("foo/bar/baz", uploadBuffer)
	if err != nil {
		t.Errorf("put failed: " + err.Error())
	}

	buffer, err := get("foo/bar/baz")
	if err != nil {
		t.Errorf("get failed: " + err.Error())
	}

	if !bytes.Equal(uploadBuffer, buffer) {
		t.Error("Storj buffers are not identical")
	}

	res, err := exists("foo/bar/baz")
	if err != nil {
		t.Errorf("exists failed: " + err.Error())
	}
	if res != true {
		t.Error("exists returned false")
	}

	err = delete("foo/bar/baz")
	if err != nil {
		t.Errorf("delete failed: " + err.Error())
	}

	res, err = exists("foo/bar/baz2")
	if err == nil {
		t.Errorf("exists failed: " + err.Error())
	}
	if res != false {
		t.Error("exists returned true")
	}
}
