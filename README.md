# bootnode-registrar

Registrar for [Geth Bootnodes](https://github.com/ethereum/go-ethereum/wiki/Setting-up-private-network-or-local-cluster#setup-bootnode). `bootnode-registrar` will resolve a DNS address record to a enode addresses that can then be consumed by `geth --bootnodes=<enodes>`. 

## Pre-requisites

`bootnode-registrar` is designed to be deployed on a Kubernetes cluster, but it *should* work anywhere. However, the following instructions make the presumption that you are on k8s:

1. Create a `bootnode` pod with 2 running containers:

    a. bootnode

        ```
        bootnode --genkey=/etc/bootnode/node.key
        bootnode --nodekey=/etc/bootnode/node.key
        ```

    b. http-server - webserver that returns the enode address

        ```
        while [ 1 ]; do echo -e \"HTTP/1.1 200 OK\n\nenode://$(bootnode -writeaddress --nodekey=/etc/bootnode/node.key)@$(MY_IP):30301\" | nc -l -v -p 8080 || break; done;
        ```

1.  Create a [headless service](https://kubernetes.io/docs/concepts/services-networking/service/#headless-services) for the `bootnode` pod(s) which will result in a DNS A record in the form of `my-namespace.svc.cluster.local`.

## Usage

```
$ git clone git@github.com:jpoon/bootnode-registrar.git
$ make
$ ./bootnode-registrar -service <dns-record-to-resolve>

-- or --

$ BOOTNODE_SERVIER=<dns-record-to-resolve> && ./bootnode-registrar
```

or skip all that and run the container

```
$ docker run jpoon/bootnode-registrar -p 9898:9898
```

Once running, a http server 

```
$ curl localhost:9898

enode://d61206cab2a77832ccd9a69e734dfe1e4a56fc4228698d696fc3799b0801398219d6e529cc33e84e5000d22865a4b42370fdfc1784d306e943553e9bfadd617e@10.244.2.43:30301, enode://b5e05a3169e7769b02c12ae457e3ed40ef1bb3443ea714d92c87ecb14c7d6eedb854aa83ab510f04813b1b0fbf49d39f781b0b38b1983a3688c25d6c5297929a@10.244.3.31:30301

```

