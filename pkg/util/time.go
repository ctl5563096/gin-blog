package util

import (
	"time"
)

func ReturnCurrentTime(Type string)  string {
	t1:=time.Now().Year()        //年
	t2:=time.Now().Month()       //月
	t3:=time.Now().Day()         //日
	t4:=time.Now().Hour()        //小时
	t5:=time.Now().Minute()      //分钟
	t6:=time.Now().Second()      //秒
	t7:=time.Now().Nanosecond()  //纳秒

	currentTimeData:=time.Date(t1,t2,t3,t4,t5,t6,t7,time.Local) //获取当前时间，返回当前时间Time

	var returnData string

	switch Type {
		case "first":
			returnData =  currentTimeData.Format("2006-01-02")
		case "second":
			returnData =  currentTimeData.Format("2006-01-02 15:04:05")
	}

	return returnData
}