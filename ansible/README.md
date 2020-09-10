# Ansible Operator

This is a simple example for ansible operator using the operator-sdk

## HowTo

This operator was created using the following command
```
operator-sdk init --domain=adawolfs.github.io --plugins=ansible
```

The speaker crd was created using the following command
```
operator-sdk create api --group hora.de.k8s --version v1 --kind Speaker --generate-role
```

Build and push operator image
```
make docker-build docker-push IMG=adawolfs/speaker-operator:ansible
```

Install Custom Resources
```
make install
```

Deploy Operator
```
make deploy IMG=adawolfs/speaker-operator:ansible
```

Create Sample Custom Resource
```
kubectl create -f ansible/config/samples/hora.de.k8s_v1_speaker.yaml 
```