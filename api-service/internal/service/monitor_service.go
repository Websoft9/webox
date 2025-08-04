package service

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type MonitorService interface {
	WriteMetrics(measurement string, tags map[string]string, fields map[string]interface{}) error
	QueryMetrics(query string) ([]map[string]interface{}, error)
	GetServerMetrics(serverID string, timeRange string) ([]map[string]interface{}, error)
	GetApplicationMetrics(appID string, timeRange string) ([]map[string]interface{}, error)
}

type monitorService struct {
	influxClient influxdb2.Client
	writeAPI     api.WriteAPI
	queryAPI     api.QueryAPI
}

func NewMonitorService(influxClient influxdb2.Client) MonitorService {
	writeAPI := influxClient.WriteAPI("websoft9", "metrics")
	queryAPI := influxClient.QueryAPI("websoft9")

	return &monitorService{
		influxClient: influxClient,
		writeAPI:     writeAPI,
		queryAPI:     queryAPI,
	}
}

func (s *monitorService) WriteMetrics(measurement string, tags map[string]string, fields map[string]interface{}) error {
	p := influxdb2.NewPoint(measurement, tags, fields, time.Now())
	s.writeAPI.WritePoint(p)
	return nil
}

func (s *monitorService) QueryMetrics(query string) ([]map[string]interface{}, error) {
	result, err := s.queryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	for result.Next() {
		record := result.Record()
		data = append(data, map[string]interface{}{
			"time":  record.Time(),
			"value": record.Value(),
			"field": record.Field(),
		})
	}

	return data, nil
}

func (s *monitorService) GetServerMetrics(serverID, timeRange string) ([]map[string]interface{}, error) {
	query := `from(bucket: "metrics")
		|> range(start: ` + timeRange + `)
		|> filter(fn: (r) => r["_measurement"] == "server_metrics")
		|> filter(fn: (r) => r["server_id"] == "` + serverID + `")`

	return s.QueryMetrics(query)
}

func (s *monitorService) GetApplicationMetrics(appID, timeRange string) ([]map[string]interface{}, error) {
	query := `from(bucket: "metrics")
		|> range(start: ` + timeRange + `)
		|> filter(fn: (r) => r["_measurement"] == "app_metrics")
		|> filter(fn: (r) => r["app_id"] == "` + appID + `")`

	return s.QueryMetrics(query)
}
