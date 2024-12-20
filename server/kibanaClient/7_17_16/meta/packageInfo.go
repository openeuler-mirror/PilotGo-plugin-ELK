/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-ELK licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Mon Jun 24 09:21:47 2024 +0800
 */
package meta

type PackageInfo_p struct {
	Name            string             `json:"name"`
	Version         string             `json:"version"`
	Title           string             `json:"title"`
	PolicyTemplates []PolicyTemplate_p `json:"policy_templates"`
	DataStreams     []DataStream_p     `json:"data_streams"`
}

type PolicyTemplate_p struct {
	Name     string                  `json:"name"`
	Inputs   []PolicyTemplateInput_p `json:"inputs"`
	Multiple bool                    `json:"multiple"`
}

type PolicyTemplateInput_p struct {
	Type string                 `json:"type"`
	Vars []map[string]interface{} `json:"vars"`
}

type DataStream_p struct {
	Type    string               `json:"type"`
	Dataset string               `json:"dataset"`
	Streams []DataStreamStream_p `json:"streams"`
	Package string               `json:"package"`
	Path    string               `json:"path"`
}

type DataStreamStream_p struct {
	Input   string                   `json:"input"`
	Vars    []map[string]interface{} `json:"vars"`
	Enabled bool                     `json:"enabled"`
}
