package tkgprotoform

import (
	config "tkgprotoform.com/protoform/pkg/config"
	files "tkgprotoform.com/protoform/pkg/tkgprotoform/files"
)

const mode = "script"

func Run() {
	conf := config.GetConf()
	files.WriteAllToLocal(conf)
}
