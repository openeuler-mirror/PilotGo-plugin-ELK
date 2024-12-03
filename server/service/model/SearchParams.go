/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-ELK licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: yajun <yajun@kylinos.cn>
 * Date: Tue Jul 23 14:47:26 2024 +0800
 */
package model

// 搜索参数
type SearchParams struct {
	Index     string  `json:"Index"`
	StartTime string  `json:"startTime"`
	EndTime   string  `json:"endTime"`
	Keyword   string  `json:"keyword"`
	Search    string  `json:"search"`
	// 其他搜索字段
}
