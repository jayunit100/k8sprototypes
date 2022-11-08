package tkgprotoform

import (
	config "tkgprotoform.com/protoform/pkg/config"
	files "tkgprotoform.com/protoform/pkg/tkgprotoform/files"
	tkg16 "tkgprotoform.com/protoform/pkg/tkgprotoform/files/tkg_1_6"
)

const mode = "script"

func Run() {
	conf := config.GetConf()
	files.WriteAllToLocal(conf, tkg16.Files())
}
