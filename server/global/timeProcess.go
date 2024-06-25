package global

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// 2024-06-24T18:57:00.000+08:00 => 18:57:00
func GetTime_HourMinSec(timestr string) (string, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		err = errors.Errorf("fail to get hour-min-sec: %s, %s **errstack**0", timestr, err.Error())
		return "", err
	}
	time_obj, err := time.ParseInLocation(time.RFC3339, timestr, loc)
	if err != nil {
		err = errors.Errorf("fail to get hour-min-sec: %s, %s **errstack**0", timestr, err.Error())
		return "", err
	}
	return fmt.Sprintf("%02d:%02d:%02d", time_obj.Hour(), time_obj.Minute(), time_obj.Second()), nil
}
