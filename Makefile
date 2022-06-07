
BIN_PATH = bin
SRC_PATH = cmd
CLIENT_NAME = gdpm-client
MASTER_NAME = gdpm-master
SLAVE_NAME = gdpm-slave

.PHONY: clean clean-key clean-binary

install: gdpm-client gdpm-master gdpm-slave

gdpm-client:
	go install github.com/gdpm/${SRC_PATH}/${CLIENT_NAME}
gdpm-master:
	go install github.com/gdpm/${SRC_PATH}/${MASTER_NAME}
gdpm-slave:
	go install github.com/gdpm/${SRC_PATH}/${SLAVE_NAME}

run-client:
	go run ${SRC_PATH}/${CLIENT_NAME}/main.go
run-master:
	go run ${SRC_PATH}/${MASTER_NAME}/*
run-slave:
	go run ${SRC_PATH}/${SLAVE_NAME}/main.go

clean-key:
	etcdctl del "" --from-key=true

clean-binary:
	rm -rf ./bin

clean: clean-key clean-binary
