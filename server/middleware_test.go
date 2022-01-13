package server

import (
	"fmt"
	"mime"
	"testing"
)

func TestValidating(t *testing.T) {
	// application/json map[charset:utf-8] <nil>
	fmt.Println(mime.ParseMediaType("application/json;charset=utf-8"))
}
