package template

/*
key: 查询模板id

value: 查询模板文本及处理函数
*/
var QueryTemplateMap map[string]QueryTemplateMeta

type SearchTemplateFunc func(string, map[string]interface{}) (interface{}, error)

type QueryTemplateMeta struct {
	Text string             `json:"text"` // 查询模板文本
	Func SearchTemplateFunc `json:"func"` // 查询模板请求处理函数
}

const (
	DSL_log_clusterhost_timeaxis_template = `{
		"script": {
		  "lang": "mustache",
		  "source": {
			"aggs": {
				"1": {
				  "date_histogram": {
					"field": "@timestamp",
					"fixed_interval": "{{fixed_interval}}",
					"time_zone": "Asia/Shanghai",
					"min_doc_count": 0
				  },
				  "aggs": {
					"1-1": {
					  "terms": {
						"field": "{{aggsfield}}",
						"order": {
						  "_count": "desc"
						},
						"size": 10,
						"min_doc_count": 0
					  }
					}
				  }
				}
			  },
			"size": "{{size}}",
			"script_fields": {},
			"stored_fields": [
			  "*"
			],
			"runtime_mappings": {},
			"query": {
				"bool": {
				  "must": [],
				  "filter": [
					{
					  "bool": {
						"should": [
						  {
							"match": {
							  "data_stream.dataset": "{{queryfield_datastream_dataset}}"
							}
						  }
						],
						"minimum_should_match": 1
					  }
					},
					{
					  "range": {
						"@timestamp": {
						  "format": "strict_date_optional_time",
						  "gte": "{{queryfield_range_gte}}",
						  "lte": "{{queryfield_range_lte}}"
						}
					  }
					}
				  ],
				  "should": [],
				  "must_not": []
				}
			}
		  },
		  "params": {
			"queryfield_datastream_dataset": "system.syslog",
			"queryfield_range_gte": "2024-06-24T10:55:36.185Z",
			"queryfield_range_lte": "2024-06-24T11:00:36.185Z",
			"aggsfield": "host.hostname",
			"size": 0,
			"fixed_interval": "10s"
		  }
		  }
	  	}
	
	}`

	DSL_log_hostprocess_timeaxis_template = `{
		"script": {
		  "lang": "mustache",
		  "source": {
			"aggs": {
				"1": {
				  "date_histogram": {
					"field": "@timestamp",
					"fixed_interval": "{{fixed_interval}}",
					"time_zone": "Asia/Shanghai",
					"min_doc_count": 0
				  },
				  "aggs": {
					"1-1": {
					  "terms": {
						"field": "{{aggsfield}}",
						"order": {
						  "_count": "desc"
						},
						"size": 10,
						"min_doc_count": 1
					  }
					}
				  }
				}
			  },
			"size": "{{size}}",
			"script_fields": {},
			"stored_fields": [
			  "*"
			],
			"runtime_mappings": {},
			"query": {
				"bool": {
				  "must": [],
				  "filter": [
					{
					  "bool": {
						"must": [
						  {
							"match": {
							  "data_stream.dataset": "{{queryfield_data_stream_dataset}}"	  
							}
						  },
						  {
							"term": {
								"host.hostname": "{{queryfield_hostname}}"
							}
						  }
						]
					  }
					},
					{
					  "range": {
						"@timestamp": {
						  "format": "strict_date_optional_time",
						  "gte": "{{queryfield_range_gte}}",
						  "lte": "{{queryfield_range_lte}}"
						}
					  }
					}
				  ],
				  "should": [],
				  "must_not": []
				}
			}
		  },
		  "params": {
			"queryfield_data_stream_dataset": "system.syslog",
			"queryfield_range_gte": "2024-06-24T10:55:36.185Z",
			"queryfield_range_lte": "2024-06-24T11:00:36.185Z",
			"aggsfield": "process.name",
			"queryfield_hostname": "wjq-pc",
			"size": 0,
			"fixed_interval": "10s"
		  }
		  }
	  	}
	
	}`

	DSL_log_stream_template = `{
		"script": {
			"lang": "mustache",
			"source": {
				"from": "{{from}}",
				"size": "{{size}}",
				"sort": [
				  {
					"@timestamp": {
					  "order": "desc",
					  "unmapped_type": "boolean"
					}
				  }
				],
				"fields": [
				  {
					"field": "*",
					"include_unmapped": "true"
				  },
				  {
					"field": "@timestamp",
					"format": "strict_date_optional_time"
				  },
				  {
					"field": "event.created",
					"format": "strict_date_optional_time"
				  },
				  {
					"field": "event.ingested",
					"format": "strict_date_optional_time"
				  }
				],
				"script_fields": {},
				"stored_fields": [
				  "*"
				],
				"runtime_mappings": {},
				"query": {
				  "bool": {
					"must": [],
					"filter": [
					  {
						"bool": {
						  "must": [
							{
							  "match": {
								"data_stream.dataset": "{{queryfield_datastream_dataset}}"
							  }
							},
							{
							  "term": {
								"host.hostname": "{{queryfield_hostname}}"
							  }
							},
							{
							  "term": {
								"process.name": "{{queryfield_processname}}"
							  }
							}
						  ]
						}
					  },
					  {
						"range": {
						  "@timestamp": {
							"format": "strict_date_optional_time",
							"gte": "{{queryfield_range_gte}}",
							"lte": "{{queryfield_range_lte}}"
						  }
						}
					  }
					],
					"should": [],
					"must_not": []
				  }
				}
			  },
			"params": {
				"queryfield_datastream_dataset": "system.syslog",
				"queryfield_range_gte": "2024-06-24T10:55:36.185Z",
				"queryfield_range_lte": "2024-06-24T11:00:36.185Z",
				"queryfield_hostname": "wjq-pc",
				"queryfield_processname": "systemd",
				"from": 0,
				"size": 10
			  }
		}
	}`
)
