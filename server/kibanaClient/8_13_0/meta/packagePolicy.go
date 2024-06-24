package meta

import "github.com/elastic/elastic-agent-libs/kibana"

type PackagePolicyInput_p struct {
	Enabled bool                             `json:"enabled"`
	Vars    map[string]interface{}           `json:"vars"`
	Streams map[string]PackagePolicyInputStream_p `json:"streams"`
}

type PackagePolicyInputStream_p struct {
	Enabled bool                   `json:"enabled"`
	Vars    map[string]interface{} `json:"vars"`
}

type PackagePolicyRequest_p struct {
	ID          string                             `json:"id,omitempty"`
	Name        string                             `json:"name"`
	Description string                             `json:"description"`
	Namespace   string                             `json:"namespace"`
	PolicyID    string                             `json:"policy_id"`
	Package     kibana.PackagePolicyRequestPackage `json:"package"`
	Vars        map[string]interface{}             `json:"vars"`
	Inputs      map[string]PackagePolicyInput_p    `json:"inputs"`
	Force       bool                               `json:"force"`
}