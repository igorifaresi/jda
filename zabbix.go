package jda

import (
	"github.com/blacked/go-zabbix"
)

const ZabbixDefaultPort = 10051

func ZabbixSendMetric(
	zabbixServerIp string,
	zabbixServerPort int,
	host string,
	itemKey string,
	value string,
) error {
	l := GetLogger()

	metrics := make([]*zabbix.Metric, 1)
    metrics[0] = zabbix.NewMetric(host, itemKey, value)

    packet := zabbix.NewPacket(metrics)

	_, err := zabbix.NewSender(zabbixServerIp, zabbixServerPort).Send(packet)
	if err != nil {
		l.Error(err.Error())
		l.Error("Error in send metric to zabbix server") //improve
		return l.ErrorQueue
	}

	return nil
}
