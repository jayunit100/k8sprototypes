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

	dir := ".generated"
	expected := []string{
		dir + "/cluster.yaml",
		dir + "/0_install_payload.sh",
		dir + "/2_management_cluster.sh",
		dir + "/2_swap_tkr.sh",
		dir + "/3_workload_cluster.sh",
		dir + "/image-builder.json",
		dir + "/image-builder-credentials.json",
		dir + "/image-builder-customizations.json",
		dir + "/image-builder-tkg.json",
		dir + "/image_builder.sh",
	}

	c := &conf.Config{}
	c.Debug = true
	if util.FileExists(dir) {
		os.RemoveAll(dir)
	}
	os.Mkdir(dir, 0777)
	c.OutputFilesPath = dir

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
