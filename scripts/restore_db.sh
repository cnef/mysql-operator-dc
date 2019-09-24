#!/usr/bin/env bash

set -o pipefail
set -o nounset
set -o xtrace

# const
RESTORE_BIN="./mysql"

# variables
__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# function
function restore_db() {
	sql=$1
	echo "${RESTORE_BIN} -h127.0.0.1 -P30446 -uroot -pUVlY88m9suHLsthK < ${sql}"
	${RESTORE_BIN} -h"127.0.0.1" -P30446 -uroot -pUVlY88m9suHLsthK < ${sql}
}

# ==> start here
if [ ! -f ${RESTORE_BIN} ]; then
	echo "${RESTORE_BIN} not exsit"
	exit 1
fi

if [ $# -lt 1 ]; then
	echo "usage: $0 db1.sql [db2.sql, db3.sql, ...]"
	exit 2
fi

for sql in "$@"; do
	restore_db "$sql"
done
