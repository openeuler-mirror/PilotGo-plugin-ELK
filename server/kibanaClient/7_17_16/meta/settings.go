package meta

type FleetOutput_p struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Is_default  bool     `json:"is_default"`
	Type        string   `json:"type"`
	Hosts       []string `json:"hosts"`
	Config_yaml string   `json:"config_yaml"`
}

type FleetOutputsResponse_p struct {
	Items []FleetOutput_p `json:"items"`
}
