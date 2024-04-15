SRC_DIR=./protocol
DST_DIR=.
PROTO_FILE=judge.proto

protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/$PROTO_FILE
# --go_opt=paths=source_relative