## Build

Variable "HELM_VERSION" and "KUBECTL_VERSION" must be passed as docker env variables during the image build i.e.

```
$ docker build --no-cache --build-arg HELM_VERSION=3.4.1 --build-arg KUBECTL_VERSION=1.19.0 -t name:tag .
```

## Notes

To completely remove Helm based application, use the following uninstall script template:

```
helm delete release_name -n namespace
kubectl delete helmchart HelmChart.metadata.name -n kube-system
kubectl delete ns namespace
```

For example:

```
helm delete kubenav -n kubenav
kubectl delete helmchart kubenav -n kube-system
kubectl delete ns kubenav
```
