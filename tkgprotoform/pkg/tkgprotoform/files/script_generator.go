package files

import (
	"fmt"
	"k8s.io/klog/v2"
	"strings"
	config "tkgprotoform.com/protoform/pkg/config"
	"tkgprotoform.com/protoform/pkg/util"
)

// PURPOSE: To write yaml and shell scripts that are custom installation artifacts to disk
// so that we can replicate the work our CI is doing locally, easily, without parsing python in our heads, or
// copying auto generated log dumps into shell scripts.
// INPUT: compile time, scans all yaml and sh files
// OUTPUT: writes the files to cwd, when you call WriteAllToLocal

import (
	_ "embed"
)

// WriteAllToLocal writes the shell scripts and yaml files
// to your local directory.  It takes a config.yaml as input.
// returns debugging info.
func WriteAllToLocal(conf *config.Config, files map[string]string) []string {
	returnVal := []string{}
	for file, contents := range files {
		outputFileLoc := func() string {
			return conf.OutputFilesPath + "/" + file
		}
		fmt.Println(outputFileLoc())
		if util.FileExists(outputFileLoc()) {
			klog.Infof("File exists %v , not writing...", file)
			returnVal = append(returnVal, fmt.Sprintf("skip %v", file))
		} else {
			klog.Infof("File not exists %v, writing...", file)

			// mostly a no-op right now
			contents = Hydrate(conf, contents)
			output := util.WriteStringToFile(contents, outputFileLoc())
			if output != nil {
				fmt.Println("ERROR ", output)
				returnVal = append(returnVal, "ERROR")
			} else {
				returnVal = append(returnVal, "success_"+outputFileLoc())
			}
		}
	}
	return returnVal
}

// Hydrate reads your configuration and returns
func Hydrate(conf *config.Config, contents string) string {

	// image builder substitutions...
	contents = strings.ReplaceAll(contents, "blah_iso_xyz", conf.ImageBuilderInputs["iso_path"])

	// read all the inputs and swap them out
	for k, v := range conf.TanzuInputs {
		contents = strings.ReplaceAll(contents, fmt.Sprintf("$%v", k), v)
	}
	// read all the inputs and swap them out
	for k, v := range conf.ImageBuilderInputs {
		contents = strings.ReplaceAll(contents, fmt.Sprintf("$%v", k), v)
	}

	// other substitutions...
	return contents
}
