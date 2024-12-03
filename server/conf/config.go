/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-ELK licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Mon Jun 24 09:21:47 2024 +0800
 */
package conf

import (
	"flag"
	"fmt"
	"os"
	"path"

	"gitee.com/openeuler/PilotGo-plugin-elk/server/global"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var Global_Config *ServerConfig

const config_type = "elk.yaml"

var config_dir string

type ServerConfig struct {
	Elk           *ElkConf
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
		fmt.Printf("%s\n", err.Error())
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
