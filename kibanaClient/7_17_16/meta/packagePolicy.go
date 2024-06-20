package meta

type PackagePolicyInput_p struct {
	Type    string                       `json:"type"`
	Enabled bool                         `json:"enabled"`
	Vars    map[string]interface{}       `json:"vars"`
	Streams []PackagePolicyInputStream_p `json:"streams"`
}

type PackagePolicyInputStream_p struct {
	Enabled     bool                                `json:"enabled"`
	Data_stream PackagePolicyInputStremDatastream_p `json:"data_stream"`
	Vars        map[string]interface{}              `json:"vars"`
}

type PackagePolicyInputStremDatastream_p struct {
	Type    string `json:"type"`
	Dataset string `json:"dataset"`
}

type PackagePolicyRequestPackage_p struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Title   string `json:"title"`
}

type PackagePolicyRequest_p struct {
	ID          string                        `json:"id,omitempty"`
	Enabled     bool                          `json:"enabled"`
	Output_id   string                        `json:"output_id"`
	Name        string                        `json:"name"`
	Description string                        `json:"description"`
	Namespace   string                        `json:"namespace"`
	PolicyID    string                        `json:"policy_id"`
	Package     PackagePolicyRequestPackage_p `json:"package"`
	Vars        map[string]interface{}        `json:"vars"`
	Inputs      []PackagePolicyInput_p        `json:"inputs"`
	Force       bool                          `json:"force"`
}
