# Golang Operator

This is a simple example for golang operator using the operator-sdk

## HowTo

This operator was created using the following command
```
operator-sdk init --domain=adawolfs.github.io --repo=github.com/adawolfs/operator-framework/golang
```

The speaker crd was created using the following command
```
operator-sdk create api --group hora.de.k8s --version v1 --kind Session --resource=true --controller=true
```

Build and push operator image

```
make docker-build docker-push IMG=adawolfs/speaker-operator:go
```

Install Custom Resources
```
make install
```

Deploy Operator
```
make deploy IMG=adawolfs/speaker-operator:go
```

Create Sample Custom Resource
```
kubectl create -f golang/config/samples/hora.de.k8s_v1_speaker.yaml 
```