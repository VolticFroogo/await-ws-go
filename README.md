# Await WS Go
[![Go Build / Test Status](https://github.com/VolticFroogo/await-ws-go/workflows/Go/badge.svg)](https://github.com/VolticFroogo/await-ws-go/actions?query=workflow%3AGo)
[![Code Coverage](https://codecov.io/gh/VolticFroogo/await-ws-go/branch/master/graph/badge.svg)](https://codecov.io/gh/VolticFroogo/await-ws-go)

Simplify WebSockets by allowing you to await upon requests.

We do this by building on top of an existing WebSocket solution, not replacing it,
so you can still keep things exactly how you want, but with these simple requests on top.

- [Await WS Go](#await-ws-go)
    - [The problem](#the-problem)
    - [The solution](#the-solution)
- [Libraries](#libraries)
- [Getting started](#getting-started)
    - [Dependencies](#dependencies)
    - [Create a client](#create-a-client)
    - [Send a request](#send-a-request)
    - [Handle a request](#handle-a-request)
    - [Handle a response](#handle-a-response)

### The problem

As WebSockets are just constant streams of data in two directions, it's impossible to
simply send a request and await a response as you would using HTTP.

This either leaves you always expecting data, or being forced to design your own
layer on top of WebSockets to allow this medium of communication.

### The solution

This library implements the latter of the two solutions, adding an additional layer on top
of WebSockets allowing you to send a request, and wait on a channel for a response.

## Libraries

Due to the nature of WebSockets, to use this library, both the server and client must use a form
of this library, or implement the protocols used in this library.

If a library is unavailable for your desired language, feel free to
[create a library request](https://github.com/VolticFroogo/await-ws-go/issues/new?labels=enhancement&template=library-request.md). 

Language | Type              | Repository
-------- | ----------------- | ----------
Go       | Server and Client | [await-ws-go](https://github.com/VolticFroogo/await-ws-go)
Dart     | Client            | [await-ws-dart](https://github.com/VolticFroogo/await-ws-dart)

## Getting started

This will be using code from [the greet example](examples/greet), so if you need to refer
to the code in larger chunks, check out that example.

Keep in mind, this is not a guide on WebSockets, just this library,
so it will assume you have knowledge on WebSockets and already have a working client.
If you don't know much about WebSockets, you can always just look at
[the greet example](examples/greet).

### Dependencies

This library only depends on [gorilla/websocket](https://github.com/gorilla/websocket),
the most widely used WebSocket library.

I decided against using [golang.org/x/net](https://godoc.org/golang.org/x/net),
the default WebSocket implementation in Golang, as it does not comply with
[RFC 6455](http://tools.ietf.org/html/rfc6455) and generally performs worse than
[gorilla/websocket](https://github.com/gorilla/websocket).

### Create a client

After establishing a WebSocket connection (whether as a server or client),
you will need to create a client.

```go
// Create a new awaitws client with the connection.
awaitClient := awaitws.NewClient(c)
```

### Send a request

Now you have a client, you can create requests which you can await upon (like an HTTP request).

As this is following [the greet example](examples/greet), we'll just send our name as the parameter.
This message can also be complex data such as structs.

```go
// Send a message to the WS server.
wait, err := awaitClient.Request("client")
if err != nil {
    log.Print(err)
    return
}
```

Now you have "wait", which is a channel, you can wait on it for a response.

```go
// Wait for a response.
response := <-wait

// Print the response message.
// Expected output: "Hello client"
log.Print(response.Message)
```

### Handle a request

Now if we can send requests, we need to be able to handle them.

Instead of decoding the incoming messages as you usually would, you'll need to decode them
as a more generic type to allow awaitws to understand it.
```go
// Decode the incoming message as a mapstructure.
var msg map[string]interface{}
err = c.ReadJSON(&msg)
```

Now you have the generic message, you can check if it's a request.
(Messages don't have to be requests following our philosophy that this library should
not be a replacement to a WebSocket library.)

Assuming the message is a request, we will just respond with it prepending "Hello "
to their message.

```go
// Verify that the message is a request.
if awaitClient.IsRequest(msg) {
    // Respond to the request prepending "Hello " to their message.
    err = awaitClient.Respond(msg, "Hello "+msg["message"].(string))
    if err != nil {
        log.Print(err)
        continue
    }
}
```

### Handle a response

Now if we can send responses, we need to be able to handle them.

NOTE: This requires you decode the message in the generic way used in [Handle a request](#handle-a-request).

This is the easiest part of implementation, you just need to make sure this function is ALWAYS
called when a message is received, so to ensure this, it's advised to call this just after
decoding the message.

If this function returns true, it means it was a response, and await-ws will handle it.
So you can just ignore this message, and it will be handled where you awaited it.

```go
// If the message is a response, await-ws will handle it, just continue.
if awaitClient.HandleResponse(msg) {
    continue
}
```