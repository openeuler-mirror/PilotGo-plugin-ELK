/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-ELK licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Tue Jun 25 10:36:07 2024 +0800
 */
package global

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// 2024-06-24T18:57:00.000+08:00 => 18:57:00
func GetTime_DateTime2HourMinSec(date_time string) (string, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		err = errors.Errorf("fail to get hour-min-sec: %s, %s **errstack**0", date_time, err.Error())
		return "", err
	}
	time_obj, err := time.ParseInLocation(time.RFC3339, date_time, loc)
	if err != nil {
		err = errors.Errorf("fail to get hour-min-sec: %s, %s **errstack**0", date_time, err.Error())
		return "", err
	}
	return fmt.Sprintf("%02d:%02d:%02d", time_obj.Hour(), time_obj.Minute(), time_obj.Second()), nil
}

func GetTime_Timestamp2DateTime(timestamp int64) string {
	t := time.Unix(timestamp/1000, 0)
	date_time := t.UTC().Format(time.RFC3339)
	return date_time
}

func GetTime_UTCDateTime2ShanghaiDateTime(utc_date_str string) (string, error) {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		err = errors.Errorf("fail to convert utc datetime to asia.shanghai datetime: %s, %s **errstack**0", utc_date_str, err.Error())
		return "", err
	}
	utctime, err := time.Parse(time.RFC3339, utc_date_str)
	if err != nil {
		err = errors.Errorf("fail to convert utc datetime to asia.shanghai datetime: %s, %s **errstack**0", utc_date_str, err.Error())
		return "", err
	}
	shanghaitime := utctime.In(location)
	return shanghaitime.Format("2006-01-02 15:04:05"), nil
}