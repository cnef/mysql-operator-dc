### 操作步骤

#### Deploy mysql cluster
```
0. modify yaml
    * cluster-with-3-replicas.yaml
        * 如果需要DR，则填写clusterDRHost字段，字段为DC端的router-service-host
        * 修改baseServerId。DC和DR的baseServerId不能重复且[baseServerId, baseServerId+members]不能重叠
        * 注意调整volumeClaimTemplate中的storage大小
    * pv*.yaml
        * 现在使用的是手动创建pv的方式，注意在nodes上预先创建spec.local.path
        * 注意修改spec.nodeAffinity.required.nodeSelectorTerm.values为pv需要安放的node节点IP
    * cluster-router.yaml
        * spec.image修改为集群可访问到的registry
    * operator/03-deployment.yaml
        * spec.template.spec.containers.image修改为mysql-operator在registry里面的imageID   

1. create pv
kubectl create -f pv0.yaml
kubectl create -f pv1.yaml
kubectl create -f pv2.yaml

2. create configmap
kubectl create configmap mycnf --from-file=my.cnf

3. create namespace
kubectl create -f 00-namespace.yaml

4. create crd
kubectl create -f 01-resources.yaml

5. create rbac
kubectl create -f 02-rbac.yaml

6. create mysql-oprator
kubectl create -f 03-deployment.yaml

7. create mysql-cluster
kubectl create -f cluster-with-3-replicas.yaml

8. create router
kubectl create -f cluster-router.yaml

9. 如果需要DR的话，则在另外个集群重复上面步骤即可
```

#### backup mysql
```
0. modify mysql-on-demand-backup.yaml
    * spec.storageProvider.s3.endpoint　minio服务端
    * spec.storageProvider.s3.bucket　存储备份的桶
    * spec.storageProvider.s3.credentialsSecret minio的验证信息

1. kubectl create -f mysql-on-demand-backup.yaml
```

#### restore mysql
```
0. modify mysql-on-demand-restore.yaml
    * spec.cluster　修改为集群名称
    * spec.backup　修改为backup的名称

1. kubectl create -f mysql-on-demand-restore.yaml
```

#### TODO
自动化脚本