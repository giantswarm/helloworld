[![](https://godoc.org/github.com/giantswarm/helloworld?status.svg)](http://godoc.org/github.com/giantswarm/helloworld)
[![](https://img.shields.io/docker/pulls/giantswarm/helloworld.svg)](https://hub.docker.com/r/giantswarm/helloworld/)
[![Go Report Card](https://goreportcard.com/badge/github.com/giantswarm/helloworld)](https://goreportcard.com/report/github.com/giantswarm/helloworld)
[![IRC Channel](https://img.shields.io/badge/irc-%23giantswarm-blue.svg)](https://kiwiirc.com/client/irc.freenode.net/#giantswarm)

# Hello World

A minimal example application for Giant Swarm.

See [The Annotaded Hello World](http://docs.giantswarm.io/guides/annotated-helloworld/) in the Giant Swarm documentation for details.

## Mac

See https://docs.giantswarm.io/reference/cli/installation/#macosx for Mac OS installation instructions.

```
$ git clone https://github.com/giantswarm/helloworld.git && cd helloworld
$ swarm login
```

(Enter your Giant Swarm credentials)

```
$ swarm up --var=domain=helloworld-$USER.gigantic.io
```

## Linux

See https://docs.giantswarm.io/reference/cli/installation/#linux for Linux installation instructions.

```
$ git clone https://github.com/giantswarm/helloworld.git && cd helloworld
$ swarm login
```

(Enter your Giant Swarm credentials)

```
$ ./swarm up --var=domain=helloworld-$USER.gigantic.io
```

## Building

Note: This project builds a simple helloworld image written in Go. The files in here are not intented for direct usage of Giant Swarm. For a suitable getting started example see http://docs.giantswarm.io/guides/your-first-service/ That said: Go around have a look. ;-)

This project is setup with an automatic build in the Docker hub: https://registry.hub.docker.com/u/giantswarm/helloworld/. To trigger an automatic build, push the changes into https://github.com/giantswarm/helloworld.
