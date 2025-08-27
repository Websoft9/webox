#!/bin/bash
set -e

if [ -f /home/appuser/data/.influxdb_initialized ]; then
    echo "influxdb already initialized, skipping..."
    exit 0
fi

: "${INFLUXDB_USERNAME:=admin}"
: "${INFLUXDB_PASSWORD:=changeme}"
: "${INFLUXDB_ORG:=myorg}"
: "${INFLUXDB_BUCKET:=mybucket}"
: "${INFLUXDB_TOKEN:=mytoken}"

# 等待 InfluxDB 启动
echo "Waiting for InfluxDB to be ready..."
until curl -s http://localhost:8086/health | grep '"status":"pass"' >/dev/null; do
  sleep 2
done

# 初始化 InfluxDB（只在首次运行时生效）
curl -XPOST http://localhost:8086/api/v2/setup \
  -H "Content-Type: application/json" \
  -d "{
        \"username\":\"${INFLUXDB_USERNAME}\",
        \"password\":\"${INFLUXDB_PASSWORD}\",
        \"org\":\"${INFLUXDB_ORG}\",
        \"bucket\":\"${INFLUXDB_BUCKET}\",
        \"token\":\"${INFLUXDB_TOKEN}\",
        \"retentionPeriodSeconds\":0
      }"

touch /home/appuser/data/.influxdb_initialized
echo "InfluxDB setup completed."