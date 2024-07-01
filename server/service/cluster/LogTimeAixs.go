package cluster

import (
	"sort"

	"gitee.com/openeuler/PilotGo-plugin-elk/server/elasticClient"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

func ProcessLogTimeAixsData(index string, querybody map[string]interface{}) (interface{}, error) {
	search_result_body_bytes, err := elasticClient.Global_elastic.SearchByTemplate(index, querybody)
	if err != nil {
		err = errors.Wrap(err, "fail to process log timeaxis data")
		return nil, err
	}

	results := []map[string]interface{}{}
	log_type_datas_map := map[string][][]int64{}
	time_agg_raw_arr := gjson.GetBytes(search_result_body_bytes, "aggregations.1.buckets").Array()

	empty_timestamp_doccount_map := InitEmptyLogTimeaxisData(time_agg_raw_arr)

	// 遍历查询结果的时间轴数组，获取每个日志类型的时间轴数组
	for _, time_agg_data_raw := range time_agg_raw_arr {
		timestamp := time_agg_data_raw.Get("key").Int()
		for _, type_doc_data_raw := range time_agg_data_raw.Get("1-1.buckets").Array() {
			time_doc_arr := []int64{}
			time_doc_arr = append(time_doc_arr, timestamp)
			time_doc_arr = append(time_doc_arr, type_doc_data_raw.Get("doc_count").Int())
			log_type_datas_map[type_doc_data_raw.Get("key").String()] = append(log_type_datas_map[type_doc_data_raw.Get("key").String()], time_doc_arr)
		}
	}

	for log_type, log_type_datas := range log_type_datas_map {
		// 补全相对空的日志时间轴数组缺失的时间戳
		for ts, empty_time_doc_arr := range empty_timestamp_doccount_map {
			for _, time_doc_arr := range log_type_datas {
				if ts == time_doc_arr[0] {
					empty_time_doc_arr[1] = time_doc_arr[1]
					break
				}
			}
		}

		// 日志时间轴数组按照时间戳升序排序
		ts_arr := []int{}
		for ts := range empty_timestamp_doccount_map {
			ts_arr = append(ts_arr, int(ts))
		}
		sort.Ints(ts_arr)
		data := [][]int64{}
		for _, ts := range ts_arr {
			data = append(data, empty_timestamp_doccount_map[int64(ts)])
		}

		results = append(results,
			map[string]interface{}{
				"name": log_type,
				"data": data,
			},
		)

		empty_timestamp_doccount_map = InitEmptyLogTimeaxisData(time_agg_raw_arr)
	}
	return results, nil
}
