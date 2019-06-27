#! /bin/bash

# build project
rm -rf ./artifacts
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o artifacts/metrics-auth
cat <<'EOF'>> ./artifacts/metric-auth.conf
{
    "redis_server": {
        "redis_host": "172.17.0.1",
        "redis_port": "6379",
        "endpoints": [
            "badger.test.net",
            "banana.test.net"
        ]
    },
    "web_server": {
        "listen_port": "8080"
    }
}
EOF

# build docker container

docker build -t morfien101/metric-auth:latest .