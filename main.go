package main

import (
    "context"
    "io"
    "log"
    "net/http"
    "time"
)

const (
    timeout         = 30 * time.Second
    maxBodySize     = 10 << 20 // 10 MB
    readHeaderTimeout = 10 * time.Second
)

type ProxyServer struct {
    client *http.Client
}

func NewProxyServer() *ProxyServer {
    return &ProxyServer{
        client: &http.Client{
            Timeout: timeout,
            Transport: &http.Transport{
                MaxIdleConns:        100,
                IdleConnTimeout:     90 * time.Second,
                DisableCompression:  true,
                MaxConnsPerHost:     10,
                DisableKeepAlives:   false,
            },
        },
    }
}

func (p *ProxyServer) ProxyHandler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), timeout)
    defer cancel()

    // Log the incoming request
    log.Printf("Received Incoming Request: %s %s", r.Method, r.URL)

    // Limit the request body size
    r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

    // Construct the target URL
    if r.URL.Host == "" {
        r.URL.Host = r.Host
    }
    if r.URL.Scheme == "" {
        r.URL.Scheme = "http"
    }
    targetURL := r.URL.String()

    // Create the outgoing request
    outgoingRequest, err := http.NewRequestWithContext(ctx, r.Method, targetURL, r.Body)
    if err != nil {
        log.Printf("Error creating request: %v", err)
        http.Error(w, "Failed to create request", http.StatusInternalServerError)
        return
    }

    // Selectively copy headers
    copyHeaders(outgoingRequest.Header, r.Header)

    // Send the request to the target server
    response, err := p.client.Do(outgoingRequest)
    if err != nil {
        log.Printf("Error sending request: %v", err)
        http.Error(w, "Failed to forward request", http.StatusBadGateway)
        return
    }
    defer response.Body.Close()

    // Copy response headers
    copyHeaders(w.Header(), response.Header)
    w.WriteHeader(response.StatusCode)

    // Copy response body with timeout
    written, err := io.Copy(w, response.Body)
    if err != nil {
        log.Printf("Error copying response: %v", err)
        return
    }
    log.Printf("Successfully proxied %d bytes", written)
}

// copyHeaders selectively copies headers from src to dst
func copyHeaders(dst, src http.Header) {
    // List of headers to forward
    allowedHeaders := map[string]bool{
        "Content-Type":     true,
        "Content-Length":   true,
        "Accept":          true,
        "Accept-Encoding": true,
        "User-Agent":      true,
        // Add other headers as needed
    }

    for key, values := range src {
        if allowedHeaders[key] {
            for _, value := range values {
                dst.Add(key, value)
            }
        }
    }
}

func main() {
    proxy := NewProxyServer()

    server := &http.Server{
        Addr:              ":8080",
        Handler:           http.HandlerFunc(proxy.ProxyHandler),
        ReadTimeout:       timeout,
        WriteTimeout:      timeout,
        ReadHeaderTimeout: readHeaderTimeout,
        MaxHeaderBytes:    1 << 20, // 1 MB
    }

    log.Printf("Starting proxy server on port %s", server.Addr)
    if err := server.ListenAndServe(); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}
