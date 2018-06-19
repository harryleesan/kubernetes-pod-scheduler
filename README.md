## Usage

In the Kubernetes Deployment manifest of the Deployments that you would like to
scale, add the annotations _scaleUp_ and
_scaleDown_.

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: app
  annotations: #used to attach data
    scaleDown: '0'
    scaleUp: '1'
spec:
  replicas: 1
  revisionHistoryLimit: 10
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 0
      maxSurge: 2
  template:
    metadata:
      labels:
        app: app
    spec:
    ...
```

```bash
docker build -t halosan/k8-pod-scheduler:latest .
docker run -v $HOME/.kube:/root/.kube -e "SCALE=scaleUp" halosan/k8-pod-scheduler:latest 
```

## Development

### Development through Docker

```bash
docker build -t halosan/k8-pod-scheduler:dev -f Dockerfile-dev .
docker run --rm -it -v $(pwd):/go/src/app \
  -v $HOME/.kube:/home/1000/.kube \
  halosan/k8-pod-scheduler:dev bash
```

### Install dependencies

From the root directory of this project:

```bash
glide up -v
```

This will create a `vendor` directory which has the correct version of Kubernetes `client-go`.

