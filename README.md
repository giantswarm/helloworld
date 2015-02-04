# Hello World

A minimal example application for Giant Swarm.

See [The Annotaded Hello World](http://docs.giantswarm.io/guides/annotated-helloworld/) in the Giant Swarm documentation for details.

## Mac

```
$ brew tap giantswarm/swarm && brew install swarm-client 
$ git clone https://github.com/giantswarm/helloworld.git && cd helloworld
$ swarm login 
```
(Enter your Giant Swarm credentials)
```
$ swarm up --var=domain=helloworld-$USER.gigantic.io
```

## Linux 

```
$ git clone https://github.com/giantswarm/helloworld.git && cd helloworld
$ curl -O http://downloads.giantswarm.io/swarm/clients/0.8.0/swarm-0.8.0-linux-amd64.tar.gz
$ tar -xzf swarm-0.8.0-linux-amd64.tar.gz swarm
$ ./swarm login 
```
(Enter your Giant Swarm credentials)
```
$ ./swarm up --var=domain=helloworld-$USER.gigantic.io
```
