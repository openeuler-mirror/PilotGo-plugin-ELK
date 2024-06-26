package cluster

import (
	"gitee.com/openeuler/PilotGo-plugin-elk/server/elasticClient"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

func ProcessLogTimeAixsData(index string, querybody map[string]interface{}) ([]map[string]interface{}, error) {
	search_result_body_bytes, err := elasticClient.Global_elastic.SearchByTemplate(index, querybody)
	if err != nil {
		err = errors.Wrap(err, "fail to process log timeaxis data")
		return nil, err
	}

	log_type_datas := []map[string]interface{}{}
	hostname_agg_raw_arr := gjson.GetBytes(search_result_body_bytes, "aggregations.1.buckets").Array()
	for _, log_type_data_raw := range hostname_agg_raw_arr {
		log_type_data := map[string]interface{}{}
		log_timestamp_datas := [][]interface{}{}

		log_type_data["name"] = log_type_data_raw.Get("key").String()
		for _, log_timestamp_data_raw := range log_type_data_raw.Get("1-1.buckets").Array() {
			log_timestamp_data := []interface{}{}
			log_timestamp_data = append(log_timestamp_data, log_timestamp_data_raw.Get("key").Int()/1000)
			log_timestamp_data = append(log_timestamp_data, log_timestamp_data_raw.Get("doc_count").Int())
			log_timestamp_datas = append(log_timestamp_datas, log_timestamp_data)
		}

		log_type_data["data"] = log_timestamp_datas
		log_type_datas = append(log_type_datas, log_type_data)
	}
	return log_type_datas, nil
}
