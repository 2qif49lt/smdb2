#!/bin/bash

DB2HOME=/home/db2inst1/sqllib
export CGO_LDFLAGS=-L$DB2HOME/lib
export CGO_CFLAGS=-I$DB2HOME/include

go build .
