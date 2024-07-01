package service

import (
	"gitee.com/openeuler/PilotGo-plugin-elk/server/global/template"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/service/cluster"
)

type SearchTemplateFunc func(string, map[string]interface{}) (interface{}, error)

var TemplateHandleFuncMap map[string]SearchTemplateFunc

func init() {
	TemplateHandleFuncMap = make(map[string]SearchTemplateFunc)
	for template_id, content := range template.DSL_template_map {
		switch content[1] {
		case "ProcessLogTimeAixsData":
			TemplateHandleFuncMap[template_id] = cluster.ProcessLogTimeAixsData
		case "ProcessLogStreamData":
			TemplateHandleFuncMap[template_id] = cluster.ProcessLogStreamData
		}
	}
}
