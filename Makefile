build:
	docker build -t pd-nginx:latest ./nginx/
	docker build -t pd-podlister  ./podlister/
tag:
	docker tag pd-podlister:latest clvx/podlister:latest
	docker tag pd-nginx:latest clvx/nginx:alpine-1.18
push:
	docker push clvx/nginx:alpine-1.18
	docker push clvx/podlister:lastest

namespace:
	kubectl create namespace pd

configure:
	#Needs configuration
	helm install stable/metrics-server --name metrics-server

install:
	helm install pd -n pd ./swarm

load:
	curl -L https://goo.gl/S1Dc3R | bash -s 500 "localhost:8080"

generate-secret:
	kubectl create secret generic podlister --from-literal=spaces-key=$(echo $SPACES_KEY) --from-literal=spaces-secret=$(echo $SPACES_SECRET) --dry-run=client -o yaml
.PHONY: .build .push .namespace .configure .load .install

