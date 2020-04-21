# UNIX Sockets based Approach

This branch is to test a new idea to base the whole concepts on top of UNIX sockets.

## Concept

1. All events are routed via buses
2. A bus is a unix sockets bound to a file
3. The binding location can be everywhere
4. Approx could therefore be used completely without a YAML file
   - Start buses ate specific locations
   - Pipe processes in and out those buses

## Code Quickstart

Use the following client server example as a quick start:

```go
// Server

package main

import (
    "log"
    "net"
)

func echoServer(c net.Conn) {
    for {
        buf := make([]byte, 512)
        nr, err := c.Read(buf)
        if err != nil {
            return
        }

        data := buf[0:nr]
        println("Server got:", string(data))
        _, err = c.Write(data)
        if err != nil {
            log.Fatal("Write: ", err)
        }
    }
}

func main() {
    l, err := net.Listen("unix", "/tmp/echo.sock")
    if err != nil {
        log.Fatal("listen error:", err)
    }

    for {
        fd, err := l.Accept()
        if err != nil {
            log.Fatal("accept error:", err)
        }

        go echoServer(fd)
    }
}
```

```go
// Client

package main

import (
    "io"
    "log"
    "net"
    "time"
)

func reader(r io.Reader) {
    buf := make([]byte, 1024)
    for {
        n, err := r.Read(buf[:])
        if err != nil {
            return
        }
        println("Client got:", string(buf[0:n]))
    }
}

func main() {
    c, err := net.Dial("unix", "/tmp/echo.sock")
    if err != nil {
        panic(err)
    }
    defer c.Close()

    go reader(c)
    for {
        _, err := c.Write([]byte("hi"))
        if err != nil {
            log.Fatal("write error:", err)
            break
        }
        time.Sleep(1e9)
    }
}
```

## Message Frames

Reliable message frames must be ensured, so that the receiver knows when it's save to attempt parsing a message. Unfortunately in many script languages a message from stdin is chunked arbitrarily by some byte size set by the used library. Several solutions are possible:

### A: Process restart for each message

Message handler is easy to implement. Performance goes down due to constant restarts of the processes.

### B: Base64 as wire-format with separator

Support in all possible scripting languages, though some implementation overhead, still easy to do.
Choose a suitable message separator string like "\n---\n" -> easy to parse, pretty-print newlines will be ignored by every base64 parser.

### B: Length-prefixed base64 wire-format

Support in all possible scripting languages, though some implementation overhead, still easy to do.
Very effective to implement with a state machine.

### Unsafe solution: Usage of Unix Process Signals

Instead of sending event-splitter via the bus, signals could be used to mark the divider between messages:
[Signals default actions](https://en.wikipedia.org/wiki/Signal_(IPC)#Default_action).

#### Possible signals:

- SIGURG
- SIGUSR2

#### Side effects

Signals are asynchronous, so that the actual end of message might be very hard to detect.
