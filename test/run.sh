#!/usr/bin/env bash

# THIS IS NOT WORKING AS A SHELLSCRIPT
# THE JOBS ARE GETTINGS DISCONNECTED AND DISAPPEARCH FROM THE JOBS LIST

./hub pipe response-pipe &
./hub fork web-server-out read-file-in find-media-type-in &
./hub merge read-file-out find-media-type-out merge-response-in &

cat response-pipe.rd | ./web-server.js > web-server-out.wr &
cat read-file-in.rd | ./read-file.js > read-file-out.wr &
cat find-media-type-in.rd | ./find-media-type.js > find-media-type-out.wr &
cat merge-response-in.rd | ./merge-response.js > response-pipe.wr &
