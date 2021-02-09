## Build

Variable "HELM_VERSION" and "KUBECTL_VERSION" must be passed as docker env variables during the image build i.e.

```
$ docker build --no-cache --build-arg HELM_VERSION=3.4.1 --build-arg KUBECTL_VERSION=1.19.0 --build-arg ACCESS_TOKEN_USR=${ACCESS_TOKEN_USR} --build-arg ACCESS_TOKEN_PWD=${ACCESS_TOKEN_PWD} -t civo/bizaar-daemon:v1alpha1 .
```
