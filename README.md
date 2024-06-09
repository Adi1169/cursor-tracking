# Cursor Tracking Application

This is a simple cursor tracking application built with Go. It uses WebSockets for real-time communication and Redis for data storage.

## Features

- Real-time cursor tracking
- WebSocket communication
- Redis data storage

## Prerequisites

- Go (version 1.16 or later recommended)
- Redis server

## Setup

1. **Clone the repository**

   Use the following command to clone this repository:

   ```bash
   git clone https://github.com/yourusername/cursor-tracking-app.git
2. **Install Go**

    Ensure Go 1.22 is [installed](https://go.dev/doc/install) and added to PATH
3. Start Redis server

    This application uses Redis for data storage. Make sure you have Redis installed and running on your machine. If not, you can download it from the [official website](https://redis.io/downloads/).

4. Run the application

    Navigate to the server directory and run the following command:
    ```bash
    go run main.go


## Test
Use the basic client provided in the repo to test the real time cursor tracking.