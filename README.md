# gdpm

gdpm is a distributed process manager written in golang. The project contains
three different software components (see the following list) and use [etcd](https://github.com/etcd-io/etcd) as the state manager.

1. gdpm-master
2. gdpm-slave
3. gdpm-client

gdpm-master is the master node that includes an API server and a scheduler. The API server listens for the http requst (ex. spawn a job or update current jobs) from the client (gdpm-client or web interface). The scheduler schedules jobs using the round-robin method. If there is no available slave node, the scheduler puts the job into a queue for later dispatching.

gdpm-slave is the worker node that is responsible for running the actual job. When gdpm-slave is online, it registers its existence with the gdpm-master and receives a unique id from the master. Then, gdpm-slave registers a listener with the etcd and listens for key changes which contain its id. When the gdpm-master scheduler schedules a job, it adds a specific format of key-value pair to etcd, so the gdpm-slave can read the key change and work accordingly.

gdpm-client works as a command-line tool for users to send requests to the master node and view the node and job's status.

## Installation
