# Remora

Remora is a framework for building dedicated server process applications and their clients.

Rather than linking core services directly into your app you link them into a server executable and have Remora manage IPC between it and your app. Since each server process only services one app, the server can stay as simple as possible.

Enhancements to the server (provided it remains backward compatible) are then available to the app without the need to re-build and re-deploy. Also, alternative server implementations can be maintained for different scenarios. This architecture is also a good fit in cases where core source code should not be made available to all developers.

Bear in mind that the IPC (currently named pipes) does come with a performance overhead (around 100-150 times slower than Golang channels). Some investigation needs to be done into whether a move to shared memory would provide a significant improvement.

Remora's IPC uses the `encoding/gob` package for convenience so there are some restrictions on the kinds of data that are marshalled. Custom marshalling can be achieved by implementing the `encoding/gob` package's `GobEncoder` and `GobDecoder` interfaces or the `encoding` package's `BinaryMarshaler` and `BinaryUnmarshaler` interfaces (this is also likely to be slightly faster since `gob` uses reflection but if you do it for this reason, make sure the performance boost warrants the extra maintenance).

Currently Remora only works on Linux since it makes use of `mkfifo` but a Windows implementation using named pipes shouldn't be too much work.

## Building a server

A Remora server is an executable that performs the following steps:
* Create a `remora/server.Server` instance, providing the name of the target application (as an absolute path or relative to the working dir)
* Open any pipes for the client and start the services that will service them
* Run the target process and wait for the result

A working example is at [test/server/main.go](test/server/main.go)

## Building a client

A Remora client is an executable that uses the `client.Client` interface to connect to pipes exposes by it's server process (the one that starts the client process) and sends/receives whatever values are required by the protocol that connects them.

A working example is at [test/client/main.go](test/client/main.go)
