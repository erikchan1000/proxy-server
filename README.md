# Go Proxy Server

This project is a simple HTTP proxy server written in Go. The server listens for incoming HTTP requests, forwards them to the target server, and returns the responses back to the client.

## Features

- Handles HTTP requests and responses.
- Forwards headers and body from the client to the destination server.
- Logs incoming requests for easy debugging.
- Lightweight and fast, leveraging Go's standard `net/http` package.

## Requirements

- Go 1.18 or later

## Installation

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/yourusername/go-proxy-server.git
   cd go-proxy-server
   ```

2. **Initialize the Project**:
   If you haven't already initialized a Go module in the project directory, run:
   ```bash
   go mod init github.com/yourusername/go-proxy-server
   ```

## Usage

1. **Run the Server**:
   Build and run the proxy server with:

   ```bash
   go run main.go
   ```

   The server listens on `http://localhost:8080` by default.

2. **Send Requests**:
   Use a tool like `curl` to send requests through the proxy:

   ```bash
   curl -x http://localhost:8080 http://example.com
   ```

3. **Change the Port**:
   Modify the `port` variable in `main.go` to change the listening port.

## Example Output

When you run the proxy and make a request, you'll see logs like this:

```plaintext
2025/01/16 12:00:00 Received request: GET http://example.com
2025/01/16 12:00:01 Forwarded response: 200 OK
```

---

Feel free to customize this project further to suit your needs. If you encounter any issues, open an issue on GitHub or reach out!
