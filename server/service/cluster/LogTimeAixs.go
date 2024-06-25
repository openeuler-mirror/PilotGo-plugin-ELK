package cluster

import (
	"github.com/tidwall/gjson"
)

func ProcessLogTimeAixsData(raw_results_bytes []byte) ([]map[string]interface{}, error) {
	log_type_datas := []map[string]interface{}{}

	hostname_agg_raw_arr := gjson.GetBytes(raw_results_bytes, "aggregations.1.buckets").Array()
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
