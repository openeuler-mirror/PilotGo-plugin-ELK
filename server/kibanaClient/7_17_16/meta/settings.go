/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-ELK licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Mon Jun 24 09:21:47 2024 +0800
 */
package meta

type FleetOutput_p struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Is_default  bool     `json:"is_default"`
	Type        string   `json:"type"`
	Hosts       []string `json:"hosts"`
	Config_yaml string   `json:"config_yaml"`
}

type FleetOutputsResponse_p struct {
	Items []FleetOutput_p `json:"items"`
}
