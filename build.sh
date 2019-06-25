#!/bin/sh

# bin=${PWD##*/}

echo Building approx itself...
rm -rf bin
mkdir -p bin/lib
go build -o bin/approx
echo ...done.

echo Building assistive processors in parent directory...

echo - approx_http_server
cd ../approx_http_server
go build -o ../approx/bin/lib/http_server

echo - approx_fork
cd ../approx_fork
go build -o ../approx/bin/lib/fork

echo - approx_merge
cd ../approx_merge
go build -o ../approx/bin/lib/merge

echo - approx_check
cd ../approx_check
go build -o ../approx/bin/lib/check
echo ...done.