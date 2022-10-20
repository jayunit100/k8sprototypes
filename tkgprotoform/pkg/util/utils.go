package util

import (
	"fmt"
	"os"
)

func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func StringToFile(s string, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString(s)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
