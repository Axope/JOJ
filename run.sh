#!/bin/bash

cmd=$1

case $cmd in
  test)
    sh swag.sh
    go run main.go
    ;;
  backend)
    go build -o app
    nohup ./app &
    ;;
  stop)
    pid=$(ps -ef | grep './app' | grep -v 'grep' | awk '{print $2}')
    if [ -n "$pid" ]; then
      echo "Stopping ./app process with PID: $pid"
      kill $pid
    else
      echo "No ./app process found"
    fi
    ;;
  *)
    echo "Invalid command. Use 'test', 'backend', or 'stop'."
    exit 1
    ;;
esac
