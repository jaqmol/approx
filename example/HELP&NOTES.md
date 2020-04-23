# Understanding piping

## Meaning of pipe characters:

- program stdout | program stdin
- program stdout > fifo/file 
- program stdin < fifo/file

## That means

- Example: `wc -l < input.txt > output.txt`
- Would be interpreted as: `(wc -l < input.txt) > output.txt`

## Expressing flow

- Goal: express data flow: `input.txt -> program -> output.txt`
- Solution: `program < input.txt > output.txt`


            |     ^
            V     |
         web-server.js <----------------o
               |                        |
               V                        |
            req-fork                    |
               |                        |
     o---------o---------o              |
     |                   |              |
     V                   V              |
read-file.js     find-media-type.js     |
     |                   |              |
     o---------o---------o              |
               |                        |
           resp-merge                   |
               |                        |
               V                        |
        merge-response.js --------------o


## Startup procedure

### 1st start the hubs

./hub pipe response-pipe &
./hub fork web-server-out read-file-in find-media-type-in &
./hub merge read-file-out find-media-type-out merge-response-in &

### 2nd start the programs

> **Note:**
> `web-server.js, find-media-type.js, merge-response.js, read-file.js`
> Need to be executable (chmod +x)

./web-server.js < response-pipe.rd > web-server-out.wr &
./find-media-type.js < find-media-type-in.rd > find-media-type-out.wr &
./merge-response.js < merge-response-in.rd > response-pipe.wr &
./read-file.js < read-file-in.rd > read-file-out.wr &
