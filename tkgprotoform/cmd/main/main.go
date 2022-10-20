package main

import (
	//	"corgon.com/corgon/pkg/tkgprotoform"
	"tkgprotoform.com/protoform/pkg/tkgprotoform"
)

const mode = "script"

func main() {
	// _ = pkg.GetConf()

	if mode == "script" {
		tkgprotoform.Run()
	}
}
