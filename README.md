# gdpm

gdpm is a distributed process manager written in golang. The project contains
three different executables (see the following list) and use [etcd](https://github.com/etcd-io/etcd) as a state manager.

1. gdpm-master
2. gdpm-slave
3. gdpm-client

gdpm-master is the master node that includes an API server and a scheduler. The API server listens for the http requst (ex. spawn a job or update current jobs) from the client (gdpm-client or web interface). The scheduler schedules jobs using the round-robin method. If there is no available slave node, the scheduler puts the job into a queue for later dispatching.

gdpm-slave is the worker node that is responsible for running the actual job. When gdpm-slave is online, it registers its existence with the gdpm-master and receives a unique id from the master. Then, gdpm-slave registers a listener with the etcd and listens for key changes which contain its id. When the gdpm-master scheduler schedules a job, it adds a specific format of key-value pair to etcd, so the gdpm-slave can read the key change and work accordingly.

gdpm-client works as a command-line tool for users to send requests to the master node and view the node and job's status.

## Environment Setup

### System Requirements

1. [Golang 1.18](https://go.dev/dl/)
2. [etcd](https://etcd.io/docs/v3.5/install/)

### Installation Guide

1. Clone the repository in to your ```$GOPATH```. If you haven't set up your go development environment, please refer to this [guide](https://go.dev/doc/gopath_code).

    ```bash
    # create the folder for the repository
    cd $GOPATH
    mkdir src/github.com
    cd src/github.com
    ```

    ```bash
    git clone git@github.com:capy-pl/gdpm.git
    ```

2. Start your etcd cluster. Please see the [official guide](https://etcd.io/docs/v3.5/dev-guide/local_cluster/).

3. Compile and install the executables. You must set the environment variables ```$GOPATH``` or ```$GOBIN``` before running the following commands. The command compiles the executables and move the executables to the folder ```$GOPATH/bin``` or ```$GOBIN``` if specified.

    ```bash
    make install
    ```

### API Usage

1. List all nodes. ```[http GET]http://localhost:8989/node/```

    ```json
    // response example
    {
        "Ids": [
            "5764e70d-a8da-4c75-aa44-09dc2dbc188b",
            "da409d07-f630-4b1d-9844-f8122f9c08d2"
        ],
        "ServiceNum": [
            0,
            0
        ],
        "Status": [
            1,
            1
        ],
        "Times": [
            "07 Jun 22 22:03 CST",
            "07 Jun 22 22:22 CST"
        ]
    } 
    ```

2. List services on a node.

    ```[http Get]http://localhost:8989/node/:nodeId/```

    ```json
    // response example
    {
        "Command": [
            "python /Users/phil/GoProjects/src/github.com/gdpm/test.py"
        ],
        "Ids": [
            "a11b8ada-1776-490f-a8a4-46c688ffb262"
        ],
        "Number": [
            4
        ]
    }
    ```

3. Create a service.

    ```[http POST]http://localhost:8989/service/```

    ```json
    // request example
    {
        "Command": "python test.py", // if you intend to execute a script, provide an absolute path
        "InstanceNum":  5, // number of process you want to spawn
    }
    ```

    ```json
    // response example
    a11b8ada-1776-490f-a8a4-46c688ffb262
    ```

4. Update a service.

    ```[http POST]http://localhost:8989/service/:serviceId/```

    ```json
    // request example
    {
        "InstanceNum": 5, // the number of instance of the service
    }
    ```

5. Delete a service.

    ```[http DELETE]http://localhost:8989/service/:serviceId/```
