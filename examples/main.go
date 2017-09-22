package main
import (
	"fmt"
        "gopkg.in/qntfy/kazaam.v3"
)

const (
	exampleJSON = "{\"Event-Name\":\"HEARTBEAT\",\"Core-UUID\":\"a475b338-99e5-11e7-8b36-3bdaa2eeeccc\",\"FreeSWITCH-Hostname\":\"ip-172-31-29-181\",\"FreeSWITCH-Switchname\":\"ip-172-31-29-181\",\"FreeSWITCH-IPv4\":\"172.31.29.181\",\"FreeSWITCH-IPv6\":\"::1\",\"Event-Date-Local\":\"2017-09-22 19:26:24\",\"Event-Date-GMT\":\"Fri, 22 Sep 2017 19:26:24 GMT\",\"Event-Date-Timestamp\":\"1506108384249057\",\"Event-Calling-File\":\"switch_core.c\",\"Event-Calling-Function\":\"send_heartbeat\",\"Event-Calling-Line-Number\":\"75\",\"Event-Sequence\":\"1398044\",\"Event-Info\":\"System Ready\",\"Up-Time\":\"0 years, 7 days, 12 hours, 10 minutes, 59 seconds, 872 milliseconds, 603 microseconds\",\"FreeSWITCH-Version\":\"1.9.0+git~20170911T194756Z~2aea0c329b~64bit\",\"Uptime-msec\":\"648659872\",\"Session-Count\":\"0\",\"Max-Sessions\":\"1000\",\"Session-Per-Sec\":\"30\",\"Session-Per-Sec-Last\":\"0\",\"Session-Per-Sec-Max\":\"26\",\"Session-Per-Sec-FiveMin\":\"0\",\"Session-Since-Startup\":\"10667\",\"Session-Peak-Max\":\"26\",\"Session-Peak-FiveMin\":\"0\",\"Idle-CPU\":\"90.166667\"}"
)

func main() {
        spec := `[{"operation": "shift","spec": {"Event": "Event-Name","UUID": "Core-UUID"}}]`
        kazaamTransform, _ := kazaam.NewKazaam(spec)
        kazaamOut, _ := kazaamTransform.TransformJSONStringToString(exampleJSON)
        fmt.Println(string(kazaamOut))
}
