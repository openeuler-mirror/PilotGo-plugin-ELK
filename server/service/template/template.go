package template

import (
	"strings"

	"gitee.com/openeuler/PilotGo-plugin-elk/server/elasticClient"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/errormanager"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/pluginclient"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/service/log"
	"github.com/pkg/errors"
)

func init() {
	QueryTemplateMap = map[string]QueryTemplateMeta{
		"log_clusterhost_timeaxis": {
			Text: DSL_log_clusterhost_timeaxis_template,
			Func: log.ProcessLogTimeAxisData,
		},
		"log_hostprocess_timeaxis": {
			Text: DSL_log_hostprocess_timeaxis_template,
			Func: log.ProcessLogTimeAxisData,
		},
		"log_stream": {
			Text: DSL_log_stream_template,
			Func: log.ProcessLogStreamData,
		},
	}
}

// 在elasticsearch中添加查询模板
func InitSearchTemplate() {
	if elasticClient.Global_elastic == nil {
		err := errors.Errorf("elasticClient is nil **errstackfatal**0") // err top
		errormanager.ErrorTransmit(pluginclient.Global_Context, err, true)
		return
	}

	for template_id, template_meta := range QueryTemplateMap {
		reqbody := strings.NewReader(template_meta.Text)
		_, err := elasticClient.Global_elastic.Client.PutScript(
			template_id,
			reqbody,
			elasticClient.Global_elastic.Client.PutScript.WithContext(elasticClient.Global_elastic.Ctx),
			elasticClient.Global_elastic.Client.PutScript.WithPretty(),
		)
		if err != nil {
			err = errors.Errorf("fail to put script: %s, %s **errstackfatal**0", template_id, err.Error()) // err top
			errormanager.ErrorTransmit(pluginclient.Global_Context, err, true)
		}
	}
}
