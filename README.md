# STUN Protocol Implementation

A Go implementation of the Session Traversal Utilities for NAT (STUN) protocol as defined in RFC 5389. This library provides both client and server implementations for NAT traversal, which is essential for peer-to-peer networking applications.

## Features

- **Full STUN Protocol Support**: Implements RFC 5389 specifications
- **Client & Server**: Both client and server implementations included
- **XOR-MAPPED-ADDRESS**: Support for the XOR-MAPPED-ADDRESS attribute
- **Structured Logging**: Comprehensive logging with configurable levels
- **Error Handling**: Robust error handling throughout the codebase
- **Easy API**: Simple and intuitive API design
- **Production Ready**: Suitable for production use with proper logging

## Installation

```bash
go get github.com/lai0xn/stun
```

## Quick Start

### Client Usage

```go
package main

import (
    "fmt"
    "github.com/lai0xn/stun"
)

func main() {
    // Create a STUN client
    client := stun.NewClient("stun.l.google.com:19302")
    
    // Send a binding request
    msg, err := client.Dial(&stun.Message{
        Header: stun.Header{
            Type: stun.BindingRequest,
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Get the XOR mapped address
    xorAddr, err := msg.GetXorAddr()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Public IP: %s:%d\n", xorAddr.IP, xorAddr.Port)
}
```

### Server Usage

```go
package main

import (
    "github.com/lai0xn/stun"
)

func main() {
    // Create a STUN server
    server := stun.NewServer(stun.ServerConfig{
        Addr: "127.0.0.1",
        Port: "3478",
        Logger: stun.NewDefaultLogger(),
    })
    
    // Start listening
    if err := server.Listen(); err != nil {
        log.Fatal(err)
    }
}
```

## Logging

The library includes a comprehensive logging system that can be configured for different environments:

### Development Logging

```go
logger := stun.NewLogger(stun.LoggerConfig{
    Level:      stun.DebugLevel,
    Format:     "text",
    ShowCaller: true,
})
```

### Production Logging

```go
logger := stun.NewLogger(stun.LoggerConfig{
    Level:      stun.InfoLevel,
    Format:     "json",
    ShowCaller: false,
})
```

### Log Levels

- `DebugLevel`: Detailed debug information
- `InfoLevel`: General information messages
- `WarnLevel`: Warning messages
- `ErrorLevel`: Error messages
- `FatalLevel`: Fatal errors (exits program)

## API Reference

### Client

#### `NewClient(addr string) *Client`
Creates a new STUN client with the specified server address.

#### `NewClientWithLogger(addr string, logger *Logger) *Client`
Creates a new STUN client with a custom logger.

#### `client.Dial(msg *Message) (*Message, error)`
Sends a STUN binding request and returns the response.

### Server

#### `NewServer(config ServerConfig) *Server`
Creates a new STUN server with the specified configuration.

#### `server.Listen() error`
Starts the server and begins listening for connections.

#### `server.Shutdown() error`
Gracefully shuts down the server.

### Message

#### `NewMessage(buff []byte) (*Message, error)`
Creates a new Message by parsing the provided byte buffer.

#### `message.GetAttr(t StunAttribute) (*Attribute, bool)`
Searches for a specific attribute type in the message.

#### `message.GetXorAddr() (*XorMappedAddr, error)`
Extracts the XOR-MAPPED-ADDRESS attribute from the message.

#### `message.Encode() []byte`
Converts the Message to its binary representation.

## Examples

See the `examples/` directory for complete working examples:

- `examples/client/client.go`: Basic client usage
- `examples/server/server.go`: Basic server usage

## Protocol Details

This implementation supports the core STUN protocol features:

- **Binding Request/Response**: Core message types for NAT discovery
- **XOR-MAPPED-ADDRESS**: Attribute containing the client's public IP
- **Transaction ID**: Unique identifier for each STUN transaction
- **Magic Cookie**: Protocol identifier (0x2112A442)
- **IPv4 Support**: Full IPv4 address handling

## Error Handling

The library provides comprehensive error handling with specific error types:

- `ErrAttrNotFound`: Attribute not found in message
- `ErrShortBuffer`: Buffer too short for reading
- `ErrInvalidCookie`: Invalid magic cookie
- `ErrShortWrite`: Incomplete write operation

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## References

- [RFC 5389 - Session Traversal Utilities for NAT (STUN)](https://tools.ietf.org/html/rfc5389)
- [STUN Protocol Overview](https://en.wikipedia.org/wiki/STUN)

## Related Projects

- [WebRTC](https://webrtc.org/): Web Real-Time Communication
- [TURN Protocol](https://tools.ietf.org/html/rfc5766): Traversal Using Relays around NAT
- [ICE Protocol](https://tools.ietf.org/html/rfc5245): Interactive Connectivity Establishment 