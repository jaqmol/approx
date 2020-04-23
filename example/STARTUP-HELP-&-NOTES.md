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

./hub pipe fifo-response-pipe &
./hub fork fifo-web-server-out fifo-read-file-in fifo-find-media-type-in &
./hub merge fifo-read-file-out fifo-find-media-type-out fifo-merge-response-in &

### 2nd start the programs

> **Note:**
> `web-server.js, find-media-type.js, merge-response.js, read-file.js`
> Need to be executable (chmod +x)

./web-server.js < fifo-response-pipe.rd > fifo-web-server-out.wr &
./find-media-type.js < fifo-find-media-type-in.rd > fifo-find-media-type-out.wr &
./merge-response.js < fifo-merge-response-in.rd > fifo-response-pipe.wr &
./read-file.js < fifo-read-file-in.rd > fifo-read-file-out.wr &
