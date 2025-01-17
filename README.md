# Go HTTP Proxy Server

A lightweight, configurable HTTP proxy server written in Go with focus on security, performance, and reliability.

## Features

- HTTP/HTTPS proxy support
- Connection pooling and keep-alive
- Configurable timeouts and size limits
- Selective header forwarding
- Request/response logging
- Context-based timeout handling
- Resource cleanup and proper error handling

## Requirements

- Go 1.21 or higher

## Installation

1. Clone the repository:

```bash
git clone https://github.com/erikchan1000/proxy-server.git
cd proxy-server
```

2. Build the application:

```bash
go build -o proxy-server
```

## Configuration

The proxy server can be configured through constants in the code:

```go
const (
    timeout           = 30 * time.Second
    maxBodySize       = 10 << 20 // 10 MB
    readHeaderTimeout = 10 * time.Second
)
```

### Available Configuration Options

- `timeout`: Maximum time for entire request/response cycle
- `maxBodySize`: Maximum size of request/response body
- `readHeaderTimeout`: Timeout for reading request headers
- `MaxIdleConns`: Maximum number of idle connections
- `MaxConnsPerHost`: Maximum number of connections per host
- `IdleConnTimeout`: How long to keep idle connections alive

## Usage

1. Start the server:

```bash
go run main.go
```

The server will start on port 8443 by default.

2. Make requests through the proxy:

```bash
# Basic GET request
curl -v "https://localhost:8443/get" -H "Host: httpbin.org"

# POST request with data
curl -v -X POST \
    -H "Content-Type: application/json" \
    -H "Host: httpbin.org" \
    -d '{"test": "data"}' \
    "http://localhost:8443/post"
```

## Testing

The repository includes several curl commands for testing different scenarios:

### Basic Tests

```bash
# Test GET request
curl -v "http://localhost:8443/get" -H "Host: httpbin.org"

# Test POST request
curl -v -X POST \
    -H "Content-Type: application/json" \
    -H "Host: httpbin.org" \
    -d '{"test": "data"}' \
    "http://localhost:8443/post"
```

# Test Script

```bash
./test.sh
```

### Load Testing

```bash
# Run multiple concurrent requests
for i in {1..10}; do
    curl -v "http://localhost:8443/get" -H "Host: httpbin.org" &
done
```

## Security Considerations

1. Header Handling

   - Only specified headers are forwarded
   - Host header is properly managed
   - Security-sensitive headers are filtered

2. Size Limits

   - Request body size is limited
   - Response size is monitored
   - Header size limits are enforced

3. Timeouts
   - Connection timeouts prevent resource exhaustion
   - Context-based cancellation is implemented
   - Idle connection timeouts are enforced

## Error Handling

The proxy server implements comprehensive error handling:

- Network errors
- Timeout errors
- Size limit violations
- Invalid request handling
- Resource cleanup

## Monitoring

The server logs important events:

- Incoming requests
- Forwarded request status
- Error conditions
- Response sizes

## Best Practices

1. Production Deployment

   - Use TLS in production
   - Implement rate limiting
   - Add authentication if needed
   - Monitor resource usage

2. Performance Optimization
   - Adjust connection pool settings
   - Monitor timeout values
   - Optimize based on usage patterns

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Known Issues

1. No built-in TLS support (in progress)
2. Limited header filtering (can be expanded)
3. Basic monitoring capabilities

## Roadmap

- [ ] Implement rate limiting
- [ ] Add metrics endpoint
- [ ] Add configuration file support
- [ ] Docker support

## Support

For issues and feature requests, please create an issue in the repository.
