* backup old-db-cluster
```shell
./backup_db.sh senseguard uums
```
* make sure that *.sql have all data in old-db-cluster
* stop old-db-cluster
* stop old-operator
* clean data in pv
* start new-operator
* start new-db-cluster
* check mysql-router that have linked mysql-cluster
* restore new-db-cluster
```shell
./restore_db.sh
```
