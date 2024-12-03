/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-ELK licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Mon Jun 24 09:21:47 2024 +0800
 */
package kibanaClient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitee.com/openeuler/PilotGo-plugin-elk/server/kibanaClient/8_13_0/meta"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/global"
	"github.com/elastic/elastic-agent-libs/kibana"
)

func (client *KibanaClient_v8) GetPackageInfo(ctx context.Context, pkgname, pkgversion string) (*meta.PackageInfo_p, error) {
	apiURL := fmt.Sprintf(meta.FleetPackageInfoAPI, pkgname, pkgversion)
	resp, err := client.Client.Connection.SendWithContext(ctx, http.MethodGet, apiURL, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error calling %s API: %w", meta.FleetPackageInfoAPI, err)
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	pinfo := Gjson_GetPackageInfo(bytes, "item.name", "item.policy_templates", "item.data_streams")
	return pinfo, nil
}

func (client *KibanaClient_v8) InstallFleetPackage(ctx context.Context, reqbody *meta.PackagePolicyRequest_p) (*kibana.PackagePolicyResponse, error) {
	reqBytes, err := json.Marshal(reqbody)
	if err != nil {
		return nil, fmt.Errorf("marshalling request json: %w", err)
	}

	apiURL := meta.FleetPackagePoliciesAPI
	resp, err := client.Client.Connection.SendWithContext(ctx, http.MethodPost, apiURL, nil, nil, bytes.NewReader(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("error calling %s API: %w", meta.FleetPackagePoliciesAPI, err)
	}
	defer resp.Body.Close()

	pkg_policy_resp := &kibana.PackagePolicyResponse{}
	err = global.ReadJSONResponse(resp, pkg_policy_resp)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	return pkg_policy_resp, nil
}
