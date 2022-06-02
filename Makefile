
BIN_PATH = bin
SRC_PATH = cmd
CLIENT_NAME = gdpm-client
MASTER_NAME = gdpm-master
SLAVE_NAME = gdpm-slave

all: client master slave

client:
	go build -o ${BIN_PATH}/${CLIENT_NAME} ${SRC_PATH}/${CLIENT_NAME}/main.go
master:
	go build -o ${BIN_PATH}/${MASTER_NAME} ${SRC_PATH}/${MASTER_NAME}/main.go
slave:
	go build -o ${BIN_PATH}/${SLAVE_NAME} ${SRC_PATH}/${SLAVE_NAME}/main.go

run-client:
	go run ${SRC_PATH}/${CLIENT_NAME}/main.go
run-master:
	go run ${SRC_PATH}/${MASTER_NAME}/main.go
run-slave:
	go run ${SRC_PATH}/${SLAVE_NAME}/main.go

clean-key:
	etcdctl del "" --from-key=true

clean-binary:
	rm -rf ./bin

clean: clean-key clean-binary
