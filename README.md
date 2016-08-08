# kubernetes-bot
Kubernetes Bot

### Usage
```
$ BOT_CHANNEL=$SLACK_CHANNEL \
  BOT_TOKEN=$SLACK_TOKEN \
  K8S_HOST=http://localhost:8080 \
  K8S_INSECURE=true \
  bin/kubebot
```

### Build
```
$ make src
$ make docker
```
