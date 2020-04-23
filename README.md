# APPROX HUB

Connect small, reactive processes in every programming language that can read and write to Stdin and Stdout in a flow graph to form applications.

- Communicating processes are also known as **Actors**
- Processes are interconnecting via a **Messaging-Bus**
- The form of programming is also known as **event-driven** or **reactive**
- The same programming pattern like **messaging queues** like **RabbitMQ/AMQ**, but at the moment restricted to a single machine

## Why?

- Easy to handle form of parallel/concurrent programming
- Easy to handle form of inter-process-communication between programs in different programming languages
- Faster than the same thing via sockets or network interface, also no port-management required
- Unix pipes are cool and everything that makes them usable for more problems is highly welcome 
- Communication between processes can be visualized

## Input, Output

Messages/events flow through a process/actor via `stdin` and `stdout`. Stream-connectors are expressed as named pipes. The builtin actors "fork" and "merge" are used to splice and reunite the event streams. The builtin actors "pipe" is used to rewire any `stdout` to any `stdin` via a name pipe.

## Simple Wire Format

Messages/events are JSON encoded, because it's pervasively supported. To distinguish messages from another, an additional envelope is used: the message is bas64-encoded with the delimiter `\n---\n`. Th reasons for this are pervasiveness of base64, ease of use and the resulting robustness.

## Usage

First the `hub` command line tool is used to set up the infrastructure of named pipes. Usually automated by a shell-script. Second the actual actors/processes are started with their Stdin and Stdout connected to the named pipes.

### `hub` CLI usage

```bash
$ ./hub
hub
Utility to build messaging systems by composing command line processes
pipe <name>
  Pipe message stream from <name>.wr to <name>.rd
fork <wr-name> <rd-name-1> <rd-name-2> <...>
  Fork message stream from wr-fifo into all provided rd-fifos
merge <wr-name-1> <wr-name-2> <...< <rd-name>
  Merge message stream from all provided wr-fifos into rd-fifo
input <name>
  Input JSON messages to stream them to <name>
cleanup <directory>
  Cleanup directory from fifos (wr & rd)
```

### Example:

Given the following flow-graph...

```ascii
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
```

...the following setup and startup script is used:

```bash
./hub pipe response-pipe &
./hub fork web-server-out read-file-in find-media-type-in &
./hub merge read-file-out find-media-type-out merge-response-in &

./web-server.js < response-pipe.rd > web-server-out.wr &
./find-media-type.js < find-media-type-in.rd > find-media-type-out.wr &
./merge-response.js < merge-response-in.rd > response-pipe.wr &
./read-file.js < read-file-in.rd > read-file-out.wr &
```

**The example folder contains a static web-server expressed as flow-graph of forward-communicating processes (in NodeJS).** The file `hub-messaging.js` contains the only API necessary in 36 LOCs, showcasing how simple support can be implemented.