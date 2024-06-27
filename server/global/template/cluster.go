package template

var (
	DSL_template_map map[string]string
)

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
						"field": "{{aggs_field}}",
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
							  "data_stream.dataset": "{{query_data_stream_dataset}}"
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
						  "gte": "{{query_range_gte}}",
						  "lte": "{{query_range_lte}}"
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
			"query_data_stream_dataset": "system.syslog",
			"query_range_gte": "2024-06-24T10:55:36.185Z",
			"query_range_lte": "2024-06-24T11:00:36.185Z",
			"aggs_field": "host.hostname",
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
						"field": "{{aggs_1-1_field}}",
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
							  "data_stream.dataset": "{{query_data_stream_dataset}}"	  
							}
						  },
						  {
							"term": {
								"host.hostname": "{{hostname}}"
							}
						  }
						]
					  }
					},
					{
					  "range": {
						"@timestamp": {
						  "format": "strict_date_optional_time",
						  "gte": "{{query_range_gte}}",
						  "lte": "{{query_range_lte}}"
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
			"query_data_stream_dataset": "system.syslog",
			"query_range_gte": "2024-06-24T10:55:36.185Z",
			"query_range_lte": "2024-06-24T11:00:36.185Z",
			"aggs_1-1_field": "process.name",
			"hostname": "wjq-pc",
			"size": 0,
			"fixed_interval": "10s"
		  }
		  }
	  	}
	
	}`
)

func init() {
	DSL_template_map = map[string]string{
		"log_clusterhost_timeaxis": DSL_log_clusterhost_timeaxis_template,
		"log_hostprocess_timeaxis": DSL_log_hostprocess_timeaxis_template,
	}
}
