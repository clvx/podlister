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
- Docker registry

Even though the above are the major dependencies, this has been developed using 
`Ubuntu 20.04`.

### Layout


    .
    ├── chart
    │   ├── Chart.yaml
    │   ├── templates
    │   │   ├── configmap.yaml
    │   │   ├── cronjob.yaml
    │   │   ├── deployment.yaml
    │   │   ├── _helpers.tpl
    │   │   ├── hpa.yaml
    │   │   ├── ingress.yaml
    │   │   ├── NOTES.txt
    │   │   ├── rolebinding.yaml
    │   │   ├── secrets.yaml
    │   │   ├── serviceaccount.yaml
    │   │   ├── service.yaml
    │   │   └── tests
    │   │       └── test-connection.yaml
    │   └── values.yaml
    ├── LICENSE
    ├── Makefile
    ├── README.md
    ├── REQUIREMENTS.md
    └── src
        ├── config
        │   ├── config.go
        │   └── config.yaml
        ├── Dockerfile
        ├── endpoint
        │   ├── endpoint.go
        │   └── endpoint_test.go
        ├── go.mod
        ├── go.sum
        ├── index.template
        └── main.go

    6 directories, 27 files

- `src/` is the ip pod discovery program written in golang using k8s and aws sdk 
to pull and upload the data. Podlister pulls endpoint information every minute. Then, 
it pushes this data to a Space bucket. Spaces is S3 compatible.
variables:
    - BUCKET_KEY:       BLOB key
    - BUCKET_SECRET:    BLOB secret
    - BUCKET_URL:       BLOB url
    - BUCKET_NAME:      BLOB name      
    - BUCKET_PRIVILEGE: public-read by default
    - TEMPLATE_NAME:    template name to be used to dump endpoint values. Defaulted to `index.template`
    - TEMPLATE_OUTPUT:  object name in the bucket. Defaulted to `index.html`
    - SERVICE_NAME:     k8s service to scan.
- `chart/` is the helm chart to deploy the `nginx` and `podlister` applications.
It defines the following:
- `nginx`: deployment + service + configmap + hpa.
- `podlister`: cronjob + secrets + configmap + serviceaccount + rolebinding.  A 
service account with proper permissions is necessary to pull information from the k8s API.

### Usage

    #Obtain a kubernetes cluster
    # Configure kubectl 

    #Create a Space Blob Storage
    # Obtain API credentials: https://www.digitalocean.com/community/tutorials/how-to-create-a-digitalocean-space-and-api-key
    # Export API credentials as SPACES_KEY abd SPACES_SECRET, or override values in MAKEFILE

    #Connect to your own docker registry
    docker login

    #Modify values in Makefile
    # Helm release
    # Docker repository

    #Configuring, building, pushing and deploying app.
    make build push configure deploy

    #Testing HPA
    #kubectl port-forward service/${SERVICE_NAME} 8080:80
    #Execute load test
    make test-load

### Load test 

- It would depend on your hpa values. Please refer to `hpa` in `chart/values.yaml`

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

### Load

- Current implementation has a request of 100m CPU per nginx container per pod. 
Each pod has only one nginx container. Depending on the node specs and node pool size(if autoscaling or fixed), the 
expected handled load would vary. HPA will trigger a new pod every time the combined resource usage hits 25% or more.

### Availability

Currently all the infrastructure layout is manual which increases the time to repair in case the cluster 
has a major misconfiguration. A better approach is having the infrastructure coded. As it provides reproducibility, 
a failover cluster could be created almost instantly or a cold one could be waiting to start rolling in case of failure. 
If the cloud provider has an outage, as the solution was developed on kubernetes, it could be replicated in another 
provider(again using IAC) and do failover by DNS.

## Observability

- As every container writes to stdout, having a process reading the standard out and forwarding the information with pod metadata would suffice to manage logging. Systems doing this approach: Loki, EFK, etc. 

## Security

- There are several way to apply security to a Kubernetes cluster. Best approach is following 
the Docker and Kubernetes CIS Benchmark. It provides the requirements to secure in a k8s cluster.
These might include:
    - Pod security policies(not running as root, not mounting host fs, etc).
    - RBAC, not mounting service account secrets, and having a service account for you application(review included chart).
    - Admission policies.
    - Building your own images(avoid pulling external or non verified images).
    - Verifying image signing.
