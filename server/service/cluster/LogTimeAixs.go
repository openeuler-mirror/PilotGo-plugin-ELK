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

	results := []map[string]interface{}{}
	log_type_datas_map := map[string][][]int64{}
	time_agg_raw_arr := gjson.GetBytes(search_result_body_bytes, "aggregations.1.buckets").Array()
	for _, time_agg_data_raw := range time_agg_raw_arr {
		timestamp := time_agg_data_raw.Get("key").Int()
		for _, host_doc_data_raw := range time_agg_data_raw.Get("1-1.buckets").Array() {
			time_doc_arr := []int64{}
			time_doc_arr = append(time_doc_arr, timestamp)
			time_doc_arr = append(time_doc_arr, host_doc_data_raw.Get("doc_count").Int())
			log_type_datas_map[host_doc_data_raw.Get("key").String()] = append(log_type_datas_map[host_doc_data_raw.Get("key").String()], time_doc_arr)
		}
	}

	for log_type, log_type_data := range log_type_datas_map {
		results = append(results,
			map[string]interface{}{
				"name": log_type,
				"data": log_type_data,
			},
		)
	}
	return results, nil
}
