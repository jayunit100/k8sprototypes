package files

import (
	"os"
	"testing"
	conf "tkgprotoform.com/protoform/pkg/config"
)

func TestFileTemplates(t *testing.T) {
	for filename, _ := range files("1.6") {
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
	c := &conf.Config{}
	c.Debug = true
	os.Remove("./.tests")
	os.Mkdir("./.tests", 0777)
	c.OutputFilesPath = "./.tests"
	out := WriteAllToLocal(c, "1.6")
	t.Logf("%v", out)
}
