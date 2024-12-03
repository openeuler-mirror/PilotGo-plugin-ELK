/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-ELK licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Thu Jul 4 16:12:28 2024 +0800
 */
package log

import "github.com/tidwall/gjson"

// 初始化空的日志时间轴数据
func InitEmptyLogTimeaxisData(time_buckets_raw_arr []gjson.Result) map[int64][]int64 {
	empty_timestamp_doccount := map[int64][]int64{}
	for _, time_agg_data_raw := range time_buckets_raw_arr {
		timestamp := time_agg_data_raw.Get("key").Int()
		time_doc_arr := []int64{}
		time_doc_arr = append(time_doc_arr, timestamp)
		time_doc_arr = append(time_doc_arr, int64(0))
		empty_timestamp_doccount[timestamp] = time_doc_arr
	}
	return empty_timestamp_doccount
}
