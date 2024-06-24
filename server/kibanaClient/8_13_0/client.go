package kibanaClient

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"gitee.com/openeuler/PilotGo-plugin-elk/conf"
	"gitee.com/openeuler/PilotGo-plugin-elk/errormanager"
	"gitee.com/openeuler/PilotGo-plugin-elk/global"
	"gitee.com/openeuler/PilotGo-plugin-elk/kibanaClient/8_13_0/meta"
	"gitee.com/openeuler/PilotGo-plugin-elk/pluginclient"
	"github.com/elastic/elastic-agent-libs/kibana"
)

var Global_kibana *KibanaClient_v8

type KibanaClient_v8 struct {
	Client *kibana.Client
	Ctx    context.Context
}

func InitKibanaClient() {
	cfg := &kibana.ClientConfig{
		Protocol: "http",
		Host:     conf.Global_Config.Kibana.Addr,
		Username: conf.Global_Config.Kibana.Username,
		Password: conf.Global_Config.Kibana.Password,
	}

	ki_client, err := kibana.NewClientWithConfig(cfg, "", "", "", "")
	if err != nil {
		err = errors.Errorf("failed to init kibana client: %+v **errstackfatal**0", err.Error()) // err top
		errormanager.ErrorTransmit(pluginclient.Global_Context, err, true)
		return
	}

	Global_kibana = &KibanaClient_v8{
		Client: ki_client,
		Ctx:    context.Background(),
	}
}

func (client *KibanaClient_v8) pkgInfo2PkgPolicyInputs(pinfo *meta.PackageInfo_p) map[string]meta.PackagePolicyInput_p {
	inputs := map[string]meta.PackagePolicyInput_p{}
	for _, policy_template_input := range pinfo.PolicyTemplates[0].Inputs {
		pkg_policy_input := meta.PackagePolicyInput_p{
			Enabled: true,
			Vars:    map[string]interface{}{},
			Streams: map[string]meta.PackagePolicyInputStream_p{},
		}
		for _, data_stream := range pinfo.DataStreams {
			if policy_template_input.Type == data_stream.Streams[0].Input {
				pkg_policy_input_stream := meta.PackagePolicyInputStream_p{
					Enabled: true,
					Vars:    map[string]interface{}{},
				}
				for _, data_stream_stream_var := range data_stream.Streams[0].Vars {
					pkg_policy_input_stream.Vars[data_stream_stream_var["name"].(string)] = data_stream_stream_var["default"]
				}
				pkg_policy_input.Streams[data_stream.Dataset] = pkg_policy_input_stream
			}
		}

		if policy_template_input.Vars != nil {
			for _, policy_template_input_var := range policy_template_input.Vars {
				pkg_policy_input.Vars[policy_template_input_var["name"].(string)] = policy_template_input_var["default"]
			}
		}

		inputs[fmt.Sprintf("%s-%s", pinfo.Name, policy_template_input.Type)] = pkg_policy_input
	}
	return inputs
}

/*
向kibana请求package info并生成package policy

input(key) => policy_templates.inputs[0].type

input(value).streams[0](key) => data_streams[0].dataset

input(key) == data_streams[0].streams[0].input
*/
func (client *KibanaClient_v8) ComposePackagePolicy(policyid, pkgname, pkgversion string) (*meta.PackagePolicyRequest_p, error) {
	pkginfo, err := client.GetPackageInfo(client.Ctx, pkgname, pkgversion)
	if err != nil {
		return nil, err
	}

	inputs := client.pkgInfo2PkgPolicyInputs(pkginfo)

	return &meta.PackagePolicyRequest_p{
		Name:      fmt.Sprintf("%s-%s", pkginfo.Name, global.GenerateRandomID(5)),
		Namespace: "",
		PolicyID:  policyid,
		Package: kibana.PackagePolicyRequestPackage{
			Name:    pkgname,
			Version: pkgversion,
		},
		Vars:   map[string]interface{}{},
		Inputs: inputs,
		Force:  true,
	}, nil
}
