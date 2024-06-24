package main

import (
	"gitee.com/openeuler/PilotGo-plugin-elk/conf"
	"gitee.com/openeuler/PilotGo-plugin-elk/db"
	"gitee.com/openeuler/PilotGo-plugin-elk/elasticClient"
	"gitee.com/openeuler/PilotGo-plugin-elk/errormanager"
	"gitee.com/openeuler/PilotGo-plugin-elk/handler"
	kibanaclient "gitee.com/openeuler/PilotGo-plugin-elk/kibanaClient/7_17_16"
	"gitee.com/openeuler/PilotGo-plugin-elk/logger"
	"gitee.com/openeuler/PilotGo-plugin-elk/pluginclient"
	"gitee.com/openeuler/PilotGo-plugin-elk/signal"
)

func main() {
	/*
		init config
	*/
	conf.InitConfig()

	/*
		init plugin client
	*/
	pluginclient.InitPluginClient()

	/*
		init error control
	*/
	errormanager.InitErrorManager()

	/*
		init web server
	*/
	handler.InitWebServer()

	/*
		init logger
	*/
	logger.InitLogger()

	/*
		init database
		neo4j mysql redis prometheus
	*/
	db.InitDB()

	/*
		init elasticsearch client
	*/
	elasticClient.InitElasticClient()

	/*
		init kibana client
	*/
	kibanaclient.InitKibanaClient()

	/*
		终止进程信号监听
	*/
	signal.SignalMonitoring()
}
