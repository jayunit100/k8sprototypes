package config

import (
	"corgon.com/corgon/pkg/tf"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"io"
	log "k8s.io/klog/v2"
	"strings"
)

func init() {
	conf = GetConf()
}

// Create private data struct to hold config.yaml options.
type Config struct {
	TanzuInputs        map[string]string
	OutputFilesPath    string
	Debug              bool
	ImageBuilderInputs map[string]string
}

var (
	conf *Config
)

func tfInputs() []string {
	// subset of parameters in tanzu framework that we care about.
	// todo, really, we should just slurp everything and use reflection ?
	tfInputs := []string{
		tf.ConfigVariableVsphereControlPlaneEndpoint,
		tf.ConfigVariableVsphereSSHAuthorizedKey,
		tf.ConfigVariableVsphereUsername,
		tf.ConfigVariableVspherePassword,
		tf.ConfigVariableVsphereCloneMode,
		tf.ConfigVariableAviEnable,
		tf.ConfigVariableAviPassword,
		tf.ConfigVariableVsphereDatacenter,
		tf.ConfigVariableVsphereDatastore,
		tf.ConfigVariableVsphereCloneMode,
		tf.ConfigVariableVsphereControlPlaneCustomVMXKeys,
		tf.ConfigVariableVsphereFolder,
		tf.ConfigVariableVsphereNetwork,
		tf.ConfigVariableVsphereTemplate,
		tf.ConfigVariableOSName,
		tf.ConfigVariableOSArch,
	}

	return tfInputs
}

func NewConfig() *Config {
	conf := &Config{}
	if conf.OutputFilesPath == "" {
		conf.OutputFilesPath = "./"
	}
	if conf.TanzuInputs == nil {
		conf.TanzuInputs = map[string]string{}
	}
	return conf
}

// GetConf reads a config.yaml file from the current directory and marshal
// into the conf config.yaml struct.
func GetConf() *Config {

	viper.AddConfigPath("./")
	viper.AddConfigPath("/Users/jayunit100/SOURCE/tkgprotoform/pkg/config")
	viper.AddConfigPath("../")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("%v", err)
	}
	conf = NewConfig()

	// loop through all user inputs and convert them to tanzu framework inputs
	// if they match a TF value (based on the config_constants.go file in tanzu framework)
	for _, v := range viper.AllKeys() {
		log.Infof("Checking input %v if it matches a tf input...", v)
		for _, vv := range tfInputs() {
			v = strings.ToUpper(v)
			vv = strings.ToUpper(vv)
			if v == vv {
				conf.TanzuInputs[v] = fmt.Sprintf("%v", viper.Get(v))
			}
			log.Infof("no match , compared %v %v ", vv, v)

		}
	}

	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config.yaml struct, %v", err)
	}

	if conf.Debug == false {
		log.Info("Disabling klog. Set debug to true if you wanna know whats going on.")
		// disable
		log.SetOutput(io.Discard)
		flags := &flag.FlagSet{}
		log.InitFlags(flags)
		flags.Set("logtostderr", "false")
	}
	return conf
}
