package conf

import (
	"flag"
	"fmt"
	"os"
	"path"

	"gitee.com/openeuler/PilotGo-plugin-elk/global"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var Global_Config *ServerConfig

const config_type = "elk.yaml"

var config_dir string

type ServerConfig struct {
	Elk           *ElkConf
	PilotGo       *PilotGoConf
	Logopts       *logger.LogOpts `yaml:"log"`
	Elasticsearch *ElasticConf
	Logstash      *LogstashConf
	Kibana        *KibanaConf
}

func ConfigFile() string {
	configfilepath := path.Join(config_dir, config_type)

	return configfilepath
}

func InitConfig() {
	flag.StringVar(&config_dir, "conf", "/opt/PilotGo/plugin/elk", "elk configuration directory")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -conf /path/to/elk.yaml(default:/opt/PilotGo/plugin/elk) \n", os.Args[0])
	}
	flag.Parse()

	bytes, err := global.FileReadBytes(ConfigFile())
	if err != nil {
		flag.Usage()
		err = errors.Wrapf(err, "open file failed: %s, %s", ConfigFile(), err.Error()) // err top
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	Global_Config = &ServerConfig{}

	err = yaml.Unmarshal(bytes, Global_Config)
	if err != nil {
		err = errors.Errorf("yaml unmarshal failed: %s", err.Error()) // err top
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
