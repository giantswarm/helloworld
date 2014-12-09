## Mac

$ brew tap giantswarm/swarm && brew install swarm-client 

$ git clone https://github.com/giantswarm/helloworld.git && cd helloworld

$ swarm login 

$ swarm up --var=domain=helloworld-$USER.gigantic.io

## Linux 


$ git clone https://github.com/giantswarm/helloworld.git && cd helloworld

$ curl -O http://downloads.giantswarm.io/swarm/clients/0.8.0/swarm-0.8.0-linux-amd64.tar.gz

$ tar -xzf swarm-0.8.0-linux-amd64.tar.gz swarm

$ ./swarm login 

$ ./swarm up --var=domain=helloworld-$USER.gigantic.io