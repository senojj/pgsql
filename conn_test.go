package pgsql

import (
	"context"
	"testing"
)

func TestConnect(t *testing.T) {
	err := connect(context.Background(), "postgres://jjones@127.0.0.1:5432")
	if err != nil {
		t.Fatal(err)
	}
	println("done...")
}
