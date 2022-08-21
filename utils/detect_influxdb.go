package utils

import (
	"context"

	"github.com/beego/beego/v2/core/logs"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func DetectInfluxDBOnline(hostname, token, org string) bool {
	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient(hostname, token)

	// Find organization
	_, err := client.OrganizationsAPI().FindOrganizationByName(context.Background(), org)
	if err != nil {
		logs.Info("InfluxDB is not exist! ", err.Error())
		return false
	}
	return true
}
