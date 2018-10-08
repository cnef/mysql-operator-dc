1. backup old-db-cluster

```
./backup_db.sh senseguard uums
```
2. ***!! make sure that *.sql have all data in old-db-cluster !!***
3. stop old-db-cluster
4. stop old-operator
5. clean data in pv
6. start new-operator
7. start new-db-cluster
8. check mysql-router that have linked mysql-cluster
9. restore new-db-cluster 

```
./restore_db.sh senseguard.sql uums.sql
```
