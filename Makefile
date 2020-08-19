KUBERNETES_NAMESPACE=pd

HELM_RELEASE=pd
CHART_NAME=swarm
CHART_PATH=./${CHART_NAME}
BLOB_KEY=$(shell echo ${SPACES_KEY}|tr -d '\n'|base64)
BLOB_SECRET=$(shell echo ${SPACES_SECRET}|tr -d '\n'|base64)
SERVICE_NAME=${HELM_RELEASE}-${CHART_NAME}

NGINX_REPOSITORY=clvx/nginx
NGINX_TAG=debug
NGINX_CONTEXT=./nginx
PODLISTER_REPOSITORY=clvx/podlister
PODLISTER_TAG=latest
PODLISTER_CONTEXT=./podlister


build-nginx:
	docker build -t ${NGINX_REPOSITORY}:${NGINX_TAG} ${NGINX_CONTEXT}

build-podlister:
	docker build -t ${PODLISTER_REPOSITORY}:${PODLISTER_TAG}  ${PODLISTER_CONTEXT}

push-nginx:
	docker push ${NGINX_REPOSITORY}:${NGINX_TAG}

push-podlister:
	docker push ${PODLISTER_REPOSITORY}:${PODLISTER_TAG}

configure-namespace:
	kubectl create namespace ${KUBERNETES_NAMESPACE}

configure-metrics:
	#Needs configuration
	helm install stable/metrics-server --name metrics-server

test-load:
	curl -L https://goo.gl/S1Dc3R | bash -s 500 "localhost:8080"

helm-deploy:
	helm upgrade --install ${HELM_RELEASE} --set secrets.key=${BLOB_KEY} --set secrets.secret=${BLOB_SECRET} --set configmap.podlister.name=${SERVICE_NAME} --set cronjob.image.tag=${PODLISTER_TAG} ${CHART_PATH}

build: build-nginx build-podlister

push: push-nginx push-podlister

configure: configure-namespace configure-metrics

.PHONY: .build .push .namespace .configure .load .helm-deploy

