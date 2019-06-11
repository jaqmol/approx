#!/bin/sh

# bin=${PWD##*/}

echo Please make sure the sources of the following processors are situated in the parent directory:
echo   - approx_http_server
echo   - approx_fork
echo   - approx_merge
echo   - approx_check

rm -rf bin
mkdir -p bin/lib

go build -o bin/approx

cd ../approx_http_server
go build -o ../approx/bin/lib/http_server

cd ../approx_fork
go build -o ../approx/bin/lib/fork

cd ../approx_merge
go build -o ../approx/bin/lib/merge

cd ../approx_check
go build -o ../approx/bin/lib/check
