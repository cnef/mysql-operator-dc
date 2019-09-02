#!/usr/bin/env bash

set -o pipefail
set -o nounset
set -o xtrace

# const
DUMP_BIN="./mysqldump"

# variables
__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# function
function dump_db() {
	db=$1
	echo "${DUMP_BIN} --databases ${db} --host=127.0.0.1 --port=30446 --user=root --password=UVlY88m9suHLsthK --set-gtid-purged=OFF > ${db}.sql"
	${DUMP_BIN} --databases ${db} --host=127.0.0.1 --port=30446 --user=root --password=UVlY88m9suHLsthK --set-gtid-purged=OFF > ${db}.sql
}

# ==> start here
if [ ! -f ${DUMP_BIN} ]; then
	echo "${DUMP_BIN} not exsit"
	exit 1
fi

if [ $# -lt 1 ]; then
	echo "usage: $0 db1 [db2, db3, ...]"
	exit 2
fi

for db in "$@"; do
	dump_db "$db"
done
