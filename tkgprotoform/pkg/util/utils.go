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

func WriteStringToFile(s string, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	l, err := f.WriteString(s)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return err
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
