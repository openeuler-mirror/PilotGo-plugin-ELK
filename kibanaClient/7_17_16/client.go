package kibanaClient

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"gitee.com/openeuler/PilotGo-plugin-elk/conf"
	"gitee.com/openeuler/PilotGo-plugin-elk/errormanager"
	"gitee.com/openeuler/PilotGo-plugin-elk/global"
	"gitee.com/openeuler/PilotGo-plugin-elk/kibanaClient/7_17_16/meta"
	"gitee.com/openeuler/PilotGo-plugin-elk/pluginclient"
	"github.com/elastic/elastic-agent-libs/kibana"
	"github.com/elastic/elastic-agent-libs/transport/httpcommon"
	"github.com/elastic/elastic-agent-libs/transport/tlscommon"
)

var Global_kibana *KibanaClient_v7

type KibanaClient_v7 struct {
	Client *kibana.Client
	Ctx    context.Context
}

func InitKibanaClient() {
	protocol := ""
	if conf.Global_Config.Elk.Https_enabled {
		protocol = "https"
	} else {
		protocol = "http"
	}
	cfg := &kibana.ClientConfig{
		Protocol: protocol,
		Host:     conf.Global_Config.Kibana.Addr,
		Username: conf.Global_Config.Kibana.Username,
		Password: conf.Global_Config.Kibana.Password,
		Transport: httpcommon.HTTPTransportSettings{
			Proxy:   httpcommon.DefaultHTTPClientProxySettings(),
			Timeout: 90 * time.Second,
			TLS: &tlscommon.Config{
				Enabled:          kibana.TRUE,
				VerificationMode: 1,
			},
		},
	}

	ki_client, err := kibana.NewClientWithConfig(cfg, "", "", "", "")
	if err != nil {
		err = errors.Errorf("fail to init kibana client: %+v **errstackfatal**0", err.Error()) // err top
		errormanager.ErrorTransmit(pluginclient.Global_Context, err, true)
		return
	}

	Global_kibana = &KibanaClient_v7{
		Client: ki_client,
		Ctx:    context.Background(),
	}
}

func (client *KibanaClient_v7) pkgInfo2PkgPolicyInputs(pinfo *meta.PackageInfo_p) []meta.PackagePolicyInput_p {
	inputs := []meta.PackagePolicyInput_p{}
	for _, policy_template_input := range pinfo.PolicyTemplates[0].Inputs {
		pkg_policy_input := meta.PackagePolicyInput_p{
			Type:    policy_template_input.Type,
			Enabled: true,
			Vars:    map[string]interface{}{},
			Streams: []meta.PackagePolicyInputStream_p{},
		}
		for _, data_stream := range pinfo.DataStreams {
			if policy_template_input.Type == data_stream.Streams[0].Input {
				pkg_policy_input_stream := meta.PackagePolicyInputStream_p{
					Enabled: true,
					Data_stream: meta.PackagePolicyInputStremDatastream_p{
						Type:    data_stream.Type,
						Dataset: data_stream.Dataset,
					},
					Vars: map[string]interface{}{},
				}
				for _, data_stream_stream_var := range data_stream.Streams[0].Vars {
					pkg_policy_input_stream.Vars[data_stream_stream_var["name"].(string)] = map[string]interface{}{
						"value": data_stream_stream_var["default"],
						"type":  data_stream_stream_var["type"],
					}
				}
				pkg_policy_input.Streams = append(pkg_policy_input.Streams, pkg_policy_input_stream)
			}
		}

		if policy_template_input.Vars != nil {
			for _, policy_template_input_var := range policy_template_input.Vars {
				pkg_policy_input.Vars[policy_template_input_var["name"].(string)] = policy_template_input_var["default"]
				pkg_policy_input.Vars[policy_template_input_var["name"].(string)] = map[string]interface{}{
					"value": policy_template_input_var["default"],
					"type":  policy_template_input_var["type"],
				}
			}
		}

		inputs = append(inputs, pkg_policy_input)
	}

	return inputs
}

/*
向kibana请求package info并生成package policy

input(key) == data_streams[0].streams[0].input
*/
func (client *KibanaClient_v7) ComposePackagePolicy(policyid, pkgname string) (*meta.PackagePolicyRequest_p, error) {
	pkginfo, err := client.GetPackageInfo(client.Ctx, pkgname)
	if err != nil {
		return nil, err
	}

	outputs, err := client.GetOutputs(client.Ctx)
	if err != nil {
		return nil, err
	}

	inputs := client.pkgInfo2PkgPolicyInputs(pkginfo)

	return &meta.PackagePolicyRequest_p{
		Enabled:   true,
		Name:      fmt.Sprintf("%s-%s", pkginfo.Name, global.GenerateRandomID(5)),
		Namespace: "default",
		PolicyID:  policyid,
		Output_id: outputs[0].Id,
		Package: meta.PackagePolicyRequestPackage_p{
			Name:    pkgname,
			Version: pkginfo.Version,
			Title:   pkginfo.Title,
		},
		Vars:   map[string]interface{}{},
		Inputs: inputs,
		Force:  true,
	}, nil
}
