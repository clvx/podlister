build:
	docker build -t pd-nginx:latest ./alpine/
push:
	docker tag pd-nginx:latest clvx/nginx:alpine-1.18
	docker push clvx/nginx:alpine-1.18

namespace:
	kubectl create namespace pd

configure:
	#Needs configuration
	helm install stable/metrics-server --name metrics-server

install:
	helm install pd -n pd ./swarm

load:
	curl -L https://goo.gl/S1Dc3R | bash -s 500 "localhost:8080"

.PHONY: .build .push .namespace .configure .load .install

