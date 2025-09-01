# 检查 Redis
if ! redis-cli -a "$WEBSOFT9_REDIS_PASSWORD" -h 127.0.0.1 ping | grep -q "PONG"; then
  echo "Redis not healthy"
  exit 1
fi

# 检查 InfluxDB
if ! curl -fsS http://127.0.0.1:8086/ping >/dev/null; then
  echo "InfluxDB not healthy"
  exit 1
fi

# 检查 api-service
if ! curl -fsS http://127.0.0.1:8080/health >/dev/null; then
  echo "api-service not healthy"
  exit 1
fi

echo "All services healthy"
exit 0
