$ brew tap giantswarm/swarm && brew install swarm-client 

$ git clone https://github.com/giantswarm/helloworld.git 

$ swarm login 

$ swarm up --var=domain=helloworld-$USER.gigantic.io
