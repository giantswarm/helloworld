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

See http://docs.giantswarm.io/reference/installation/#linux for Linux installation instructions.

```
$ git clone https://github.com/giantswarm/helloworld.git && cd helloworld
$ ./swarm login 
```
(Enter your Giant Swarm credentials)
```
$ ./swarm up --var=domain=helloworld-$USER.gigantic.io
```
