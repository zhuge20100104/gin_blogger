package util

import (
	"os"
	"testing"
)

func TestGetRootDir(t *testing.T) {
	t.Log(GetRootDir())
	t.Log(os.Args[0])
}
