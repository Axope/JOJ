#!/bin/bash

# MongoDB 连接信息
MONGO_HOST="localhost"
MONGO_PORT="27017"
MONGO_USER=$MONGO_USER
MONGO_PASSWORD=$MONGO_PASSWORD
DATABASE="JOJ"
COLLECTION="problems"

# 生成100个随机数据文件
for ((i=1; i<=100; i++))
do
    cat <<EOF > problem_$i.json
{
  "title": "Title $i",
  "timeLimit": $((RANDOM % 10 + 1)),
  "memoryLimit": $((RANDOM % 449 + 64)),
  "description": "Description for problem $i",
  "testSamples": [
    {
      "input": "Input sample $i",
      "output": "Output sample $i"
    }
  ],
  "dataRange": "Data range for problem $i",
  "point": $((RANDOM % 100 + 1)),
  "tags": ["Tag1", "Tag2"],
  "tutorial": "Tutorial for problem $i",
  "testCases": [
    {
      "input": "Test case input $i",
      "output": "Test case output $i"
    }
  ]
}
EOF
done

# 插入数据
for ((i=1; i<=100; i++))
do
    mongoimport --host $MONGO_HOST --port $MONGO_PORT -u $MONGO_USER -p $MONGO_PASSWORD --authenticationDatabase admin --db $DATABASE --collection $COLLECTION --file problem_$i.json
done

# 清理临时文件
rm problem_*.json

echo "100 documents inserted into 'problems' collection."
