package files

import (
	"fmt"
	"os"
	"testing"
	conf "tkgprotoform.com/protoform/pkg/config"
	tkg16 "tkgprotoform.com/protoform/pkg/tkgprotoform/files/tkg_1_6"
	"tkgprotoform.com/protoform/pkg/util"
)

func TestFileTemplates(t *testing.T) {
	for filename, _ := range tkg16.Files() {
		fileinfo, err := os.Stat(filename)
		if err != nil {
			t.Log("Missing or error file", fileinfo, err)
			t.Fail()
		} else {
			t.Log("file", fileinfo)
		}
	}
}

func TestFilesWriting(t *testing.T) {

	expected := []string{
		".tests/2_management_cluster.sh",
		".tests/cluster.yaml",
		".tests/image_builder.sh",
	}

	c := &conf.Config{}
	c.Debug = true
	os.RemoveAll("./.tests")
	os.Mkdir("./.tests", 0777)
	c.OutputFilesPath = "./.tests"

	// confirm files dont exist
	for _, v := range expected {
		if util.FileExists(v) {
			t.Fatal("dir not clean", v)
		}
	}

	out := WriteAllToLocal(c, tkg16.Files())
	t.Logf(fmt.Sprintf("%v", out))

	// confirm files DO exist
	for _, v := range expected {
		if !util.FileExists(v) {
			t.Log("missing file", v)
			t.Fail()
		}
	}
}
