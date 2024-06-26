package kibanaClient

import (
	"gitee.com/openeuler/PilotGo-plugin-elk/server/kibanaClient/8_13_0/meta"
	"github.com/tidwall/gjson"
)

/*
@varssrc: vars字段值的json数据

将varssrc转换为[]map[string]interface{}数据
*/
func VarsProcess(varssrc gjson.Result) []map[string]interface{} {
	varsdst := []map[string]interface{}{}
	varssrc.ForEach(func(key gjson.Result, value gjson.Result) bool {
		_var := map[string]interface{}{}
		for k, v := range value.Map() {
			if v.IsArray() {
				_arr := []string{}
				for _, s := range v.Array() {
					_arr = append(_arr, s.String())
				}
				_var[k] = _arr
			} else if v.IsBool() {
				_var[k] = v.Bool()
			} else {
				_var[k] = v.String()
			}
		}
		varsdst = append(varsdst, _var)
		return true
	})
	return varsdst
}

func Gjson_GetPackageInfo(pkginfobytes []byte, searchitems ...string) *meta.PackageInfo_p {
	pinfo := &meta.PackageInfo_p{
		Name:            "",
		PolicyTemplates: []meta.PolicyTemplate_p{},
		DataStreams:     []meta.DataStream_p{},
	}
	results := gjson.GetManyBytes(pkginfobytes, searchitems...)
	pinfo.Name = results[0].String()
	policy_templates_raw_arr := results[1].Array()
	policy_template := meta.PolicyTemplate_p{
		Name:     policy_templates_raw_arr[0].Get("name").String(),
		Inputs:   []meta.PolicyTemplateInput_p{},
		Multiple: policy_templates_raw_arr[0].Get("multiple").Bool(),
	}
	policy_templates_inputs_raw_arr := policy_templates_raw_arr[0].Get("inputs").Array()
	for _, policy_templates_input_raw := range policy_templates_inputs_raw_arr {
		policy_template_input := meta.PolicyTemplateInput_p{
			Type: policy_templates_input_raw.Get("type").String(),
			Vars: []map[string]interface{}{},
		}
		_vars := VarsProcess(policy_templates_input_raw.Get("vars"))
		policy_template_input.Vars = _vars
		policy_template.Inputs = append(policy_template.Inputs, policy_template_input)
	}
	pinfo.PolicyTemplates = append(pinfo.PolicyTemplates, policy_template)

	data_streams_raw_arr := results[2].Array()
	for _, data_stream_raw := range data_streams_raw_arr {
		data_stream := meta.DataStream_p{
			Type:    data_stream_raw.Get("type").String(),
			Dataset: data_stream_raw.Get("dataset").String(),
			Streams: []meta.DataStreamStream_p{},
			Package: data_stream_raw.Get("package").String(),
			Path:    data_stream_raw.Get("path").String(),
		}
		data_stream_streams_raw_arr := data_stream_raw.Get("streams").Array()
		for _, data_stream_stream_raw := range data_stream_streams_raw_arr {
			data_stream_stream := meta.DataStreamStream_p{
				Input:   data_stream_stream_raw.Get("input").String(),
				Vars:    []map[string]interface{}{},
				Enabled: data_stream_stream_raw.Get("enabled").Bool(),
			}
			_vars := VarsProcess(data_stream_stream_raw.Get("vars"))
			data_stream_stream.Vars = _vars
			data_stream.Streams = append(data_stream.Streams, data_stream_stream)
		}
		pinfo.DataStreams = append(pinfo.DataStreams, data_stream)
	}
	return pinfo
}
