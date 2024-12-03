/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-ELK licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Mon Jun 24 09:21:47 2024 +0800
 */
package conf

type ElkConf struct {
	Https_enabled      bool   `yaml:"https_enabled"`
	Public_certificate string `yaml:"cert_file"`
	Private_key        string `yaml:"key_file"`
	Addr               string `yaml:"server_listen_addr"`
	Addr_target        string `yaml:"server_target_addr"`
}

type ElasticConf struct {
	Https_enabled bool   `yaml:"https_enabled"`
	Addr          string `yaml:"addr"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
}

type LogstashConf struct {
	Addr string `yaml:"http_addr"`
}

type KibanaConf struct {
	Https_enabled bool   `yaml:"https_enabled"`
	Addr          string `yaml:"addr"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
}
