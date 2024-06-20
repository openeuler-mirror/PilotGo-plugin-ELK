package meta

const (
	FleetPackageInfoAPI     = "/api/fleet/epm/packages/%s/%s"
	FleetPackagePoliciesAPI = "/api/fleet/package_policies"
)

type PackageInfo_p struct {
	Name            string             `json:"name"`
	PolicyTemplates []PolicyTemplate_p `json:"policy_templates"`
	DataStreams     []DataStream_p     `json:"data_streams"`
}

type PolicyTemplate_p struct {
	Name     string                  `json:"name"`
	Inputs   []PolicyTemplateInput_p `json:"inputs"`
	Multiple bool                    `json:"multiple"`
}

type PolicyTemplateInput_p struct {
	Type string                   `json:"type"`
	Vars []map[string]interface{} `json:"vars"`
}

type DataStream_p struct {
	Type    string               `json:"type"`
	Dataset string               `json:"dataset"`
	Streams []DataStreamStream_p `json:"streams"`
	Package string               `json:"package"`
	Path    string               `json:"path"`
}

type DataStreamStream_p struct {
	Input   string                   `json:"input"`
	Vars    []map[string]interface{} `json:"vars"`
	Enabled bool                     `json:"enabled"`
}