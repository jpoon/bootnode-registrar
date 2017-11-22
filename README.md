# bootnode-registrar

Registrar for [Geth Bootnodes](https://github.com/ethereum/go-ethereum/wiki/Setting-up-private-network-or-local-cluster#setup-bootnode). Given a DNS A record, `bootnode-registrar` will resolve the record to a set of IPs, and will iteratively call each IP to retrieve the enode address.

## Requirements

* Kubernetes cluster

## Setup

* Create a pod to run Geth `bootnode`
    ```
    bootnode --genkey=/etc/bootnode/node.key
    bootnode --nodekey=/etc/bootnode/node.key
    ```
* In a sidecar pod, run a webserver to return the enode address
    ```
    while [ 1 ]; do echo -e \"HTTP/1.1 200 OK\n\nenode://$(bootnode -writeaddress --nodekey=/etc/bootnode/node.key)@$(MY_IP):30301\" | nc -l -v -p 8080 || break; done;
    ```
* Create a [Kubernetes Headless Service](https://kubernetes.io/docs/concepts/services-networking/service/#headless-services) which will create a DNS A Record in the form of `my-svc.my-namespace.svc.cluster.local`.

## Usage

```
$ docker run jpoon/bootnode-registrar -p 9898:9898
$ curl localhost:9898

[
"enode://d61206cab2a77832ccd9a69e734dfe1e4a56fc4228698d696fc3799b0801398219d6e529cc33e84e5000d22865a4b42370fdfc1784d306e943553e9bfadd617e@10.244.2.43:30301",
"enode://b5e05a3169e7769b02c12ae457e3ed40ef1bb3443ea714d92c87ecb14c7d6eedb854aa83ab510f04813b1b0fbf49d39f781b0b38b1983a3688c25d6c5297929a@10.244.3.31:30301"
]

```

