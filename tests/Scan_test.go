package tests

import (
	"github.com/saichler/syncit/files"
	"testing"
)

func TestScan(t *testing.T) {
	root:=files.Scan("/home/saichler/robot")
	files.Print(root,2,true,true)
	files.Stat(root,"/tmp/robot")
}
