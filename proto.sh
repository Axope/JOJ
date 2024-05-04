SRC_DIR=./protocol
DST_DIR=.
PROTO_FILE_judge=judge.proto
PROTO_FILE_judge_result=judge_result.proto

protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/$PROTO_FILE_judge
protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/$PROTO_FILE_judge_result