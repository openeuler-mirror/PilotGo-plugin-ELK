/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-ELK licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Mon Jun 24 09:21:47 2024 +0800
 */
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
	"gitee.com/openeuler/PilotGo-plugin-elk/server/service/template"
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
		init search template
	*/
	elasticClient.InitElasticClient()
	template.InitSearchTemplate()

	/*
		init kibana client
	*/
	kibanaclient.InitKibanaClient()

	/*
		终止进程信号监听
	*/
	signal.SignalMonitoring()
}
