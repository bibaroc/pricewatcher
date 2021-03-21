# pricewatcher

## how to build
```sh
docker build -f config/docker/watcher/Dockerfile -t yornesek/pricewatcher:$(git rev-parse --short HEAD) -t yornesek/pricewatcher:latest .
docker push yornesek/pricewatcher -a
```
## how to deploy
```sh
kubectl -n pricewatcher rollout restart deploy watcher
```