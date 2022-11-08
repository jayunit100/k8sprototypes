package util

import (
	"bufio"
	"fmt"
	"k8s.io/klog/v2"
	"os"
	"strings"
)

func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func SearchFileForString(filepath string, search string) bool {
	file, err := os.Open(filepath)

	if err != nil {
		klog.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()
	all := ""
	for _, eachline := range txtlines {
		all = all + eachline
	}
	if strings.Contains(all, search) {
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
