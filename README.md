## Build

### Local/Development

Variable "HELM_VERSION" and "KUBECTL_VERSION" must be passed as docker env variables during the image build i.e.

```
$ docker build --no-cache --build-arg HELM_VERSION=v3.4.1 --build-arg KUBECTL_VERSION=v1.19.0 --build-arg ACCESS_TOKEN_USR=${ACCESS_TOKEN_USR} --build-arg ACCESS_TOKEN_PWD=${ACCESS_TOKEN_PWD} -t civo/bizaar-daemon:v1alpha1 .
```

### Release

The GitHub Actions will do its job to fetch the latest kubectl and helm versions and pass them as `--build-arg` when building the image. Refer _docker-build-push.yaml_ file for more details.
