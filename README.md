# gdpm

gdpm is a distributed process manager written in golang. The project contains
three different executables (see the following list) and use [etcd](https://github.com/etcd-io/etcd) as a state manager.

## Introduction

gdpm has three software components (see the following list) and uses etcd to synchronize services states between different software components. We will give a brief introduction to each component in the following sections.

1. gdpm-master
2. gdpm-slave
3. gdpm-client

### gdpm-master

gdpm-master is the master node that includes a TCP server, an HTTP server, and a scheduler. The TCP server listens for the registration request from gdpm-slaves. When a gdpm-slave joins the cluster, it registers its existence with the master so that the master can assign services to it. The HTTP server listens for an http request (ex. spawn a new service or update current services) from the client (gdpm-client or web interface). If the request changes a service's status, the master will update the key-value pair in etcd accordingly. The scheduler selects a slave from the registered slave pool and schedules the service using the round-robin method. If there is no available slave, the scheduler puts the service into a queue for later dispatching.

### gdpm-slave

gdpm-slave is the worker node that actually runs a service. When gdpm-slave is online, it registers its existence with the gdpm-master and receives a unique id from the master. Then, gdpm-slave registers a listener with the etcd and listens for key changes which contain its id. When the gdpm-master scheduler schedules a job, it adds a specific key-value pair format to etcd, so the gdpm-slave can read the key change and work accordingly.

### gdpm-client

gdpm-client works as a command-line tool for users to send requests to the master node and view the node and job's status.

## Environment Setup

### System Requirements

1. [Golang 1.18](https://go.dev/dl/)
2. [etcd](https://etcd.io/docs/v3.5/install/)
3. [node](https://nodejs.org/en/download/) is required if you need the web GUI.

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

3. Install the dependencies, compile and install the executables. You must set the environment variables ```$GOPATH``` or ```$GOBIN``` before running the following commands. The command compiles the executables and move the executables to the folder ```$GOPATH/bin``` or ```$GOBIN``` if specified.

    ```bash
    go mod download
    make install
    ```

#### How to Build and Start the Web GUI

1. Change the working directory into the ```web``` directory.

    ```bash
    cd $GOPATH/src/github.com/gdpm/web
    ```

2. Install dependencies.

    ```bash
    npm install
    ```

3. Launch development server.

    ```bash
    npm run dev
    ```

4. Open ```http://localhost:3000``` in your browser.

### Example Usage

1. Check the etcd's status.

    ```bash
    etcdctl member list
    ```

    If the etcd cluster is running, you should see the following log.

    ```bash
    8e9e05c52164694d, started, default, http://localhost:2380, http://localhost:2379, false
    ```

2. Start the master node.

    ```bash
    gdpm-master
    ```

    ```bash
    # example output
    2022/06/11 01:24:30 http server is listening at 0.0.0.0:8989
    2022/06/11 01:24:30 Server starts listening for connection at 0.0.0.0:8888.
    ```

3. Start a worker node(s).

    ```bash
    gdpm-slave
    ```

4. Review the status of the cluster. The command lists the id, status and number of services of a node.

    ```bash
    # the command list current worker nodes in the cluster
    gdpm-client get nodes
    ```

5. Start a service. The following command tells the master to schedule a service consist of n instance and each of the instance runs the given command.

    ```bash
    gdpm-client -r <instance number> <command>

    # the following command spawn 5 processes and all run the giving
    # python command
    gdpm-client -r 5 python /absolute/directory/path/test.py
    ```

6. Review the status of a worker node. The command lists all services' status on the given node.

    ```bash
    gdpm-client get node <id>
    ```

7. Besides CLI, you can also use the Web GUI to manage node and services.

### API Usage

1. List all nodes. ```[http GET]http://localhost:8989/node/```

    ```jsonc
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

    ```jsonc
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

    ```jsonc
    // request example
    {
        "Command": "python test.py", // if you intend to execute a script, provide an absolute path
        "InstanceNum":  5, // number of process you want to spawn
    }
    ```

    ```jsonc
    // response example
    a11b8ada-1776-490f-a8a4-46c688ffb262
    ```

4. Update a service.

    ```[http POST]http://localhost:8989/service/:serviceId/```

    ```jsonc
    // request example
    {
        "InstanceNum": 5, // the number of instance of the service
    }
    ```

5. Delete a service.

    ```[http DELETE]http://localhost:8989/service/:serviceId/```
