#!/bin/bash
set -e

if [ -f /home/appuser/data/.influxdb_initialized ]; then
    echo "influxdb already initialized, skipping..."
    exit 0
fi

: "${WEBSOFT9_INFLUXDB_USERNAME:=admin}"
: "${WEBSOFT9_INFLUXDB_PASSWORD:=changeme}"
: "${WEBSOFT9_INFLUXDB_ORG:=websoft9}"
: "${WEBSOFT9_INFLUXDB_BUCKET:=metrics}"
: "${WEBSOFT9_INFLUXDB_TOKEN:=mytoken}"

# 等待 InfluxDB 启动
echo "Waiting for InfluxDB to be ready..."
until curl -s http://localhost:8086/health | grep '"status":"pass"' >/dev/null; do
  sleep 2
done

# 初始化 InfluxDB（只在首次运行时生效）
curl -XPOST http://localhost:8086/api/v2/setup \
  -H "Content-Type: application/json" \
  -d "{
        \"username\":\"${WEBSOFT9_INFLUXDB_USERNAME}\",
        \"password\":\"${WEBSOFT9_INFLUXDB_PASSWORD}\",
        \"org\":\"${WEBSOFT9_INFLUXDB_ORG}\",
        \"bucket\":\"${WEBSOFT9_INFLUXDB_BUCKET}\",
        \"token\":\"${WEBSOFT9_INFLUXDB_TOKEN}\",
        \"retentionPeriodSeconds\":0
      }"

touch /home/appuser/data/.influxdb_initialized
echo "InfluxDB setup completed."