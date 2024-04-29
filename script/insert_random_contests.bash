#!/bin/bash

# MongoDB 连接信息
MONGO_HOST="localhost"
MONGO_PORT="27017"
MONGO_USER=$MONGO_USER
MONGO_PASSWORD=$MONGO_PASSWORD
DATABASE="JOJ"
COLLECTION="contests"

function random_pid() {
    local PIDS=("65fd7c63f9781e8e749a0996" "65fd7c6318bafe25d18ac71c" "65fd7c636c669cdc7cd491c8" "65fd7c64334969f08825b65f" "65fd7c64b8fc4fa9da7c39b3")
    local idx=$((RANDOM % ${#PIDS[@]}))
    echo "${PIDS[idx]}"
}
function random_status() {
    local STATUSES=("Register" "Running" "Close")
    local idx=$((RANDOM % ${#STATUSES[@]}))
    echo "${STATUSES[idx]}"
}

# 生成100个随机数据文件
for ((i=1; i<=100; i++))
do
    cat <<EOF > contest_$i.json
{
  "title": "Contest $i",
  "status": "$(random_status)",
  "startTime": "$(date -Ins)",
  "duration": $((RANDOM % 1000 + 1)),
  "note": "Random note $i",
  "problems": [
    {
      "pid": "$(random_pid)",
      "nick": "$(printf \\$(printf '%03o' $(($RANDOM % 26 + 65))))",
      "title": "Problem $i-A"
    },
    {
      "pid": "$(random_pid)",
      "nick": "$(printf \\$(printf '%03o' $(($RANDOM % 26 + 65))))",
      "title": "Problem $i-B"
    }
  ]
}
EOF
done

# 插入数据
for ((i=1; i<=100; i++))
do
    mongoimport --host $MONGO_HOST --port $MONGO_PORT -u $MONGO_USER -p $MONGO_PASSWORD --authenticationDatabase admin --db $DATABASE --collection $COLLECTION --file contest_$i.json
done

# 清理临时文件
rm contest_*.json

echo "100 documents inserted into 'contests' collection."
