# Terminal Chat Application in Go

A simple terminal-based chat application written in Go that allows multiple users to communicate in real-time through a client-server architecture.

## Features

- Real-time messaging between multiple clients
- Clean terminal interface with message history
- Nickname-based user identification
- Automatic notifications when users join or leave
- Nickname validation:
  - No empty nicknames allowed
  - No duplicate nicknames (case-insensitive)
- Clear message display with automatic screen refresh

## How It Works

The application consists of two main components:
- A server that handles multiple client connections and broadcasts messages
- A client that connects to the server and provides a user interface for sending/receiving messages

## Getting Started

### Prerequisites
- Go 1.x or higher

### Running the Server

1. Navigate to the server directory
2. Run the server:
```bash
go run server.go
```

The server will start listening on port 8080.

### Running a Client

1. Navigate to the client directory
2. Run the client:
```bash
go run client.go
```
3. Enter your nickname when prompted
4. Start chatting!


## Technical Details

- Uses TCP for network communication
- Implements concurrent client handling with goroutines
- Employs mutex for thread-safe client list management
- Uses ANSI escape sequences for terminal manipulation
