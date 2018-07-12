### 操作步骤
1. create pv
```
kubectl create -f pv0.yaml
kubectl create -f pv1.yaml
kubectl create -f pv2.yaml
```

2. create namespace
```
kubectl create -f 00-namespace.yaml
```

3. create crd
```
kubectl create -f 01-resources.yaml
```

4. create rbac
```
kubectl create -f 02-rbac.yaml
```

5. create mysql-oprator
```
kubectl create -f 03-deployment.yaml
```

6. create mysql-cluster
```
kubectl create -f cluster-with-3-replicas.yaml
```

7. create router
```
kubectl create -f cluster-router.yaml
```