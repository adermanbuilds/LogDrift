#!/bin/bash
# Generate sample log data for testing

echo "\033[1;36mGenerating sample log stream...\033[0m"
echo ""

# Simulate various log patterns
cat << 'EOF'
2024-01-03 10:23:45 INFO Application started successfully
2024-01-03 10:23:46 INFO Database connection established
2024-01-03 10:23:47 INFO Server listening on port 8080
2024-01-03 10:24:12 INFO Request received: GET /api/users
2024-01-03 10:24:13 INFO Response sent: 200 OK (45ms)
2024-01-03 10:24:45 WARN Slow query detected: SELECT * FROM users took 823ms
2024-01-03 10:25:01 ERROR Database connection timeout after 30s
2024-01-03 10:25:02 ERROR Failed to fetch user data: connection refused
2024-01-03 10:25:03 ERROR HTTP 502 Bad Gateway from upstream server
2024-01-03 10:25:05 INFO Retrying connection...
2024-01-03 10:25:06 INFO Connection re-established
2024-01-03 10:25:30 INFO Request: POST /api/orders
2024-01-03 10:25:31 WARN Rate limit approaching for IP 192.168.1.100
2024-01-03 10:25:45 ERROR Payment gateway timeout: transaction_id=12345
2024-01-03 10:26:00 FATAL Out of memory: killed process 4532
2024-01-03 10:26:01 ERROR Panic in goroutine: nil pointer dereference
2024-01-03 10:26:02 ERROR Stack trace: main.go:245
2024-01-03 10:26:03 INFO Application restarted by watchdog
2024-01-03 10:26:15 INFO Health check: OK
2024-01-03 10:26:20 WARN Disk usage at 85%
2024-01-03 10:26:45 ERROR Failed to write to log file: disk full
2024-01-03 10:27:00 INFO Request: GET /api/health
2024-01-03 10:27:01 INFO Response: 200 OK (5ms)
2024-01-03 10:27:30 ERROR Authentication failed for user: invalid_token
2024-01-03 10:27:31 WARN Suspicious activity detected from IP 10.0.0.55
2024-01-03 10:27:32 ERROR Access denied: IP blocked
2024-01-03 10:28:00 INFO Request: GET /api/dashboard
2024-01-03 10:28:01 INFO Response: 200 OK (120ms)
2024-01-03 10:28:15 WARN Cache miss: high load on database
2024-01-03 10:28:30 ERROR Database deadlock detected on table orders
2024-01-03 10:28:31 ERROR Transaction rolled back: timeout after 60s
2024-01-03 10:29:00 INFO Background job started: cleanup_old_logs
2024-01-03 10:29:15 INFO Deleted 1,523 old log entries
2024-01-03 10:29:30 INFO Background job completed in 30s
2024-01-03 10:30:00 INFO System status: All services operational
EOF

echo ""
echo "\033[1;32mTest complete. Run this script and pipe to logdrift:\033[0m"
echo "  ./test_logs.sh | ./logdrift"
