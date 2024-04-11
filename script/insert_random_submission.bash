#!/bin/bash

# MongoDB 连接信息
MONGO_HOST="localhost"
MONGO_PORT="27017"
MONGO_USER=$MONGO_USER
MONGO_PASSWORD=$MONGO_PASSWORD
DATABASE="JOJ"
COLLECTION="submissions"

function random_pid() {
    local PIDS=("65fd7c63f9781e8e749a0996" "65fd7c6318bafe25d18ac71c" "65fd7c636c669cdc7cd491c8" "65fd7c64334969f08825b65f" "65fd7c64b8fc4fa9da7c39b3")
    local idx=$((RANDOM % ${#PIDS[@]}))
    echo "${PIDS[idx]}"
}

function random_lang() {
    local LANGS=("Cpp" "Java" "Python" "Go")
    local idx=$((RANDOM % ${#LANGS[@]}))
    echo "${LANGS[idx]}"
}

function random_status() {
    local STATUSES=("Pending" "Compiling" "Judging" "Compile Error" "Accept" "Wrong Answer" "Time Limit Exceeded" "Memory Limit Exceeded" "Runtime Error" "Output Limit Exceeded" "Unknown Error")
    local idx=$((RANDOM % ${#STATUSES[@]}))
    echo "${STATUSES[idx]}"
}

function random_code() {
    cat /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c 100
}

# 生成100个随机数据文件
for ((i=1; i<=100; i++))
do
    cat <<EOF > submission_$i.json
{
  "uid": $((RANDOM % 1000 + 1)),
  "pid": { "\$oid" : "$(random_pid)" },
  "submitTime": "$(date -Ins)",
  "lang": "$(random_lang)",
  "status": "$(random_status)",
  "runningTime": $((RANDOM % 1000 + 1)),
  "runningMemory": $((RANDOM % 1024 + 1)),
  "submitCode": "$(random_code)"
}
EOF
done

# 插入数据
for ((i=1; i<=100; i++))
do
    mongoimport --host $MONGO_HOST --port $MONGO_PORT -u $MONGO_USER -p $MONGO_PASSWORD --authenticationDatabase admin --db $DATABASE --collection $COLLECTION --file submission_$i.json
done

# 清理临时文件
rm submission_*.json

echo "100 documents inserted into 'submissions' collection."
