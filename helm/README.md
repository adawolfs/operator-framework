# Helm Operator

This is a simple example for helm operator using the operator-sdk

## HowTo

This operator was created using the following command
```
operator-sdk init --plugins=helm --domain=adawolfs.github.io
```

The speaker crd was created using the following command
```
operator-sdk create api --group hora.de.k8s --version v1 --kind Speaker
```

Build and push helm operator image
```
make docker-build docker-push IMG=adawolfs/speaker-operator:helm
```

Install Custom Resources
```
make install
```

Deploy Operator
```
make deploy IMG=adawolfs/speaker-operator:helm
```

Create Sample Custom Resource
```
kubectl create -f helm/config/samples/hora.de.k8s_v1_speaker.yaml 
```