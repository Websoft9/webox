package database

import (
	"api-service/internal/config"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// InitInfluxDB initializes InfluxDB connection (unchanged)
func InitInfluxDB(cfg *config.Config) (influxdb2.Client, error) {
	client := influxdb2.NewClient(cfg.InfluxDB.URL, cfg.InfluxDB.Token)
	return client, nil
}
