# Deploying in Kubernetes

## Adding Kubernetes Secrets
```sh
kubectl create secret generic pricewatcher --from-env-file=./config/k8s/pricewatcher/secrets.pricewatcher.env -n pricewatcher
```