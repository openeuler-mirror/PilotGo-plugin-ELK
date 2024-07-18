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
