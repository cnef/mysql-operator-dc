### 操作步骤
* create pv
```
kubectl create -f pv0.yaml
kubectl create -f pv1.yaml
kubectl create -f pv2.yaml
```

* create namespace
```
kubectl create -f 00-namespace.yaml
```

* create crd
```
kubectl create -f 01-resources.yaml
```

* create rbac
```
kubectl create -f 02-rbac.yaml
```

* create mysql-oprator
```
kubectl create -f 03-deployment.yaml
```

* create mysql-cluster
```
kubectl create -f cluster-with-3-replicas.yaml
```

* create router
```
kubectl create -f cluster-router.yaml
```