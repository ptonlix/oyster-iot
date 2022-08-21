package services

import (
	"context"
	"time"

	"github.com/beego/beego/v2/core/logs"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxServer struct {
	client      influxdb2.Client
	writeClient api.WriteAPIBlocking
	readClient  api.QueryAPI
}

func NewInfluxService(hostname, token, org, bucket string) *InfluxServer {
	influxS := &InfluxServer{}
	influxS.client = influxdb2.NewClient(hostname, token)
	influxS.writeClient = influxS.client.WriteAPIBlocking(org, bucket)
	influxS.readClient = influxS.client.QueryAPI(org)

	return influxS
}

func (i *InfluxServer) WriteData(measurement string, tag map[string]string, fields map[string]interface{}) error {
	p := influxdb2.NewPoint(measurement, tag, fields, time.Now())
	// write point immediately
	if err := i.writeClient.WritePoint(context.Background(), p); err != nil {
		return err
	}
	return nil
}

func (i *InfluxServer) ReadData(query string) (*[]map[string]interface{}, error) {
	data, err := i.query(query)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (i *InfluxServer) query(query string) (*[]map[string]interface{}, error) {
	data := []map[string]interface{}{}
	// Get parser flux query result
	result, err := i.readClient.Query(context.Background(), query)
	if err == nil {
		// Use Next() to iterate over query result lines
		for result.Next() {
			// read result
			//logs.Debug("row: %#v\n", result.Record().Values())
			data = append(data, result.Record().Values())
		}
		if result.Err() != nil {
			logs.Debug("Query error: %s\n", result.Err().Error())
		}
	} else {
		return nil, err
	}
	//logs.Debug("data is %#v", data)
	return &data, nil
}

func (i *InfluxServer) queryRaw(query string) {
	// Query and get complete result as a string
	// Use default dialect
	result, err := i.readClient.QueryRaw(context.Background(), query, api.DefaultDialect())
	if err == nil {
		logs.Debug("QueryResult:")
		logs.Debug(result)
	} else {
		logs.Debug("Query error: %s\n", err.Error())
	}

}

func (i *InfluxServer) queryRawWithParams(query string) {

}

func (i *InfluxServer) Close() {
	// Ensures background processes finishes
	i.client.Close()
}
