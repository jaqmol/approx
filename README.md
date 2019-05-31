# APPROX

Application baseplane. Build applications from communicating processes. Mix and match libraries and languages as you wish.

## Nomenclature

### Hub

Reads a formation and spins up a network of processors. The hub is the `approx` executable.

### Processor

A process classified by:

- Is parametrized by environments variables
- Listens for JSON-RPC messages from one or many input-streams
- Processes commands and data provided by JSON-RPC messages
- Writes JSON-RPC messages to one or many output-streams
- A processor is supposed to be as low complexity as possible
- A procesor is expected to listen for JSON-RPC input messages until it receives SIGINT

### Input, Output

A pipe or named pipe.
If a processor has only one input and one output, messages flow through it via `stdin` and `stdout`.
If a processor has multiple inputs, they are provided by the hub as named pipes. The paths to named pipes are provided by environemt variables prefixed with "IN_" and ending on an index. For instance: "IN_0", "IN_1", or "OUT_0", "OUT_1". The total number of named pipes is provided by the environment variables "IN_COUNT", "OUT_COUNT".

## Run example

### static_file_server

File server build as a flow-chart of stream-processors.

``` bash
go build
PATH=example ENDPOINT=/ ./approx example/static_file_server/formation.json
```

## Why

- To allow for designing applications as a network or flow-chart of stream processors.
- To design each of them as a low complexity stream process that transforms an input into an output.
- To build multi-process-applications easily and with every programming language that can read/write files.
- To mix and match programming languages as comes handy.
- To choose the best library (or programming language) for a single problem. 
- Not to be forced into compromises.

### Why pipes and named pipes?

A network of stream processors is an event driven architecture. Listening for messages by reading from `stdin` or a named pipe is an equally easy to to use API for input-output in the very most of programming languages. Also:
- Very good performance.
- Easy access.
- Maximum operating system support.

### Why not unix sockets?

Too difficult to reach in most programming languages. Too complex to use.
Sockets don't perform better than pipes.

### Why not http, websockets, long polling, ...?

Request-response-based network communication protocols are not a natural or effective choice for event driven architectures.
Communication of processes via network interface is not as ressource efficient as pipes and named pipes.