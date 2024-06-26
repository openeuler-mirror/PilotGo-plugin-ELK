package main

import (
	"gitee.com/openeuler/PilotGo-plugin-elk/server/conf"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/db"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/elasticClient"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/errormanager"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/handler"
	kibanaclient "gitee.com/openeuler/PilotGo-plugin-elk/server/kibanaClient/7_17_16"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/logger"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/pluginclient"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/signal"
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
