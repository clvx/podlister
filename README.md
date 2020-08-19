## Overview

Construct a system which uses a load balancer to balance traffic between any number 
of nginx containers, which serve a static file containing the IP addresses of all 
of the currently connected nginx containers. This should all work on the fly as 
nginx containers are added or removed.

Review `REQUIREMENTS.md` for constraints details.

## Podlister

Podlister uses a cronjob to pull service endpoint information of a custom nginx 
container deployment.

### Dependencies

- Helm 3
- Kubernetes 1.17
- Kubectl v 1.17
- Go version 1.14.7
- Docker 19.03.12
- Digital Ocean Kubernetes
- Digital Ocean Spaces Storage Object

Even though the above are the major dependencies, this has been developed using 
`Ubuntu 20.04`.

### Layout
    .
    ├── Makefile
    ├── nginx
    │   ├── docker-entrypoint.sh
    │   ├── Dockerfile
    │   └── LICENSE
    ├── podlister
    │   ├── config.yaml
    │   ├── Dockerfile
    │   ├── go.mod
    │   ├── go.sum
    │   ├── index.template
    │   └── main.go
    ├── README.md
    └── swarm
        ├── Chart.yaml
        ├── templates
        │   ├── configmap.yaml
        │   ├── cronjob.yaml
        │   ├── deployment.yaml
        │   ├── _helpers.tpl
        │   ├── hpa.yaml
        │   ├── ingress.yaml
        │   ├── NOTES.txt
        │   ├── rolebinding.yaml
        │   ├── secrets.yaml
        │   ├── serviceaccount.yaml
        │   ├── service.yaml
        │   └── tests
        │       └── test-connection.yaml
        └── values.yaml

    5 directories, 25 files

- `podlister/` is the ip pod discovery program written in golang using k8s and aws sdk 
to pull and upload the data. Podlister pulls endpoint information every minute. Then, 
it pushes this data to a Space bucket. Spaces is S3 compatible.
variables:
    - BUCKET_KEY:       BLOB key
    - BUCKET_SECRET:    BLOB secret
    - BUCKET_URL:       BLOB url
    - BUCKET_NAME      
    - BUCKET_PRIVILEGE: public-read by default
    - TEMPLATE_NAME:    template name to be used to dump endpoint values. Defaulted to `index.template`
    - TEMPLATE_OUTPUT:  object name in the bucket. Defaulted to `index.html`
    - SERVICE_NAME:     k8s service to scan.
- `nginx/` is a fork of `nginx:alpine` using `alpine 3.10` without some helper scripts.
It servers the Spaces bucket information.
- `swarm/` is the helm chart to deploy the `nginx` and `podlister` applications.
It defines the following:
- `nginx`: deployment + service + configmap + hpa.
- `podlister`: cronjob + secrets + configmap + serviceaccount + rolebinding.  A 
service account with proper permissions is necessary to pull information from the k8s API.

### Usage

    #Configuring, building, pushing and deploying app.
    make configure build push helm-deploy

    #Testing HPA
    helm get notes ${RELEASE_NAME}
    #Copy the output commands which have the following structures:
    #export POD_NAME=$(kubectl ...)
    #kubectl port-forward $POD_NAME 8080:80
    #Execute load test
    make test-load

Load test would depend on your hpa values. Please refer to `hpa` in `swarm/values.yaml`

## Discussion

### Node information

The current implementation using CronJobs will have the following issues:
- It runs only over a period of time. If real time updates are necessary, then a 
deployment might be a better option with a polling logic which is always pushing the current state. 
The downside of this approach is always hitting the blob storage endpoint. In that case, using 
a distributed storage to place a file which is mounted to all nodes would be a better approach; however,
this would be limited by cluster or region boundaries.
- Another point to take in account is that the current implementation does not keep previous data; hence, 
information is for the specific moment in time. This can be fixed if files are uploaded with a unique 
string; then, another logic in the frontend can pull the right information checking the latest timestamp.
Another approach is pulling data from the bucket, and add new data until reaches a threshold then partition 
the data.
