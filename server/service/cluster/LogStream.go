package cluster

import (
	"encoding/json"

	"gitee.com/openeuler/PilotGo-plugin-elk/server/elasticClient"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/global"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

func ProcessLogStreamData(index string, querybody map[string]interface{}) (interface{}, error) {
	search_result_body_bytes, err := elasticClient.Global_elastic.SearchByTemplate(index, querybody)
	if err != nil {
		err = errors.Wrap(err, "fail to process log timeaxis data")
		return nil, err
	}

	data := map[string]interface{}{}
	returned_logs := []map[string]interface{}{}
	hits_raw_arr := gjson.GetBytes(search_result_body_bytes, "hits.hits").Array()
	for _, hit_raw := range hits_raw_arr {
		hit_map := map[string][]interface{}{}
		json.Unmarshal([]byte(hit_raw.Get("fields").Raw), &hit_map)

		log := map[string]interface{}{}
		if log["date"], err = global.GetTime_UTCDateTime2ShanghaiDateTime(hit_map["@timestamp"][0].(string)); err != nil {
			err = errors.Wrap(err, "fail to process log timeaxis data")
			return nil, err
		}
		if hit_map["log.level"] != nil {
			log["level"] = hit_map["log.level"][0].(string)
		} else {
			log["level"] = ""
		}
		log["processname"] = hit_map["process.name"][0].(string)
		log["message"] = hit_map["message"][0].(string)

		returned_logs = append(returned_logs, log)
	}

	data["total"] = gjson.GetBytes(search_result_body_bytes, "hits.total.value").Int()
	data["hits"] = returned_logs
	return data, nil
}
