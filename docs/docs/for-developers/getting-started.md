---
sidebar_position: 1
---
# Getting Started Locally

## Prerequestites
1. Install the following utilities
    - kubectl
    - kubectx (optional)
    - helm
    - skaffold

## Make sure your Kubernetes cluster is app and running
## Configure Cluster
1. Install ingress controller

```sh
helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace
```
2. Make sure ingress controller is up and running by
```sh
curl http://kubernetes.docker.internal/
```