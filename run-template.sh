#!/bin/bash

DB2HOME=/home/db2inst1/sqllib
export LD_LIBRARY_PATH=$DB2HOME/lib

./smdb2 -db=xxx -host=127.0.0.1 -port=11111 -user=xxx -pwd=xxxxxx -tarips=111.111.11.1,222.222.22.2