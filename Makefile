HELM_RELEASE=bc
CHART_NAME=swarm
CHART_PATH=./${CHART_NAME}
BUCKET_KEY=$(shell echo ${SPACES_KEY}|tr -d '\n'|base64)
BUCKET_SECRET=$(shell echo ${SPACES_SECRET}|tr -d '\n'|base64)
BUCKET_NAME=${HELM_RELEASE}-${CHART_NAME}
BUCKET_URL=https://${HELM_RELEASE}-${CHART_NAME}.nyc3.digitaloceanspaces.com/
SERVICE_NAME=${HELM_RELEASE}-${CHART_NAME}

NGINX_REPOSITORY=clvx/nginx
NGINX_TAG=alpine-1.18
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

helm-metrics:
	helm repo add stable https://kubernetes-charts.storage.googleapis.com/
	helm upgrade metrics-server stable/metrics-server --install --namespace kube-system --set args={--kubelet-preferred-address-types=InternalIP}

test-load:
	curl -L https://goo.gl/S1Dc3R | bash -s 100 "localhost:8080"

helm-deploy:
	helm upgrade --install ${HELM_RELEASE}  --set secrets.key=${BUCKET_KEY} --set secrets.secret=${BUCKET_SECRET} --set configmap.podlister.serviceName=${SERVICE_NAME} --set configmap.podlister.bucketName=${BUCKET_NAME} --set cronjob.image.tag=${PODLISTER_TAG} --set configmap.nginx.proxy=${BUCKET_URL} ${CHART_PATH}

build: build-nginx build-podlister

push: push-nginx push-podlister

configure: helm-metrics

deploy: helm-deploy

.PHONY: .build .push .configure .test-load .deploy

