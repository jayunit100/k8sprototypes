package files

import (
	"os"
	"testing"
	conf "tkgprotoform.com/protoform/pkg/config"
)

func TestFiles(t *testing.T) {

	c := &conf.Config{}

	os.Remove("./.tests")
	os.Mkdir("./.tests", 0777)
	c.OutputFilesPath = "./.tests/"
	WriteAllToLocal(c)

	for filename, _ := range files() {
		fileinfo, err := os.Stat(filename)
		if err != nil {
			t.Log("Missing or error file", fileinfo, err)
			t.Fail()

		} else {
			t.Log("file", fileinfo)
		}
	}
}
