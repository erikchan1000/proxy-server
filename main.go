package main

import (
    "context"
    "crypto/tls"
    "io"
    "log"
    "net/http"
    "time"
)

const (
    timeout           = 30 * time.Second
    maxBodySize       = 10 << 20 // 10 MB
    readHeaderTimeout = 10 * time.Second
)

type ProxyServer struct {
    client *http.Client
}

func NewProxyServer() *ProxyServer {
    // Create custom transport with TLS config
    transport := &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true, // Only for testing
        },
        MaxIdleConns:        100,
        IdleConnTimeout:     90 * time.Second,
        DisableCompression:  true,
        MaxConnsPerHost:     10,
        DisableKeepAlives:   false,
    }

    return &ProxyServer{
        client: &http.Client{
            Timeout:   timeout,
            Transport: transport,
        },
    }
}

func (p *ProxyServer) ProxyHandler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), timeout)
    defer cancel()

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

    // Copy headers
    copyHeaders(outgoingRequest.Header, r.Header)

    // Send the request
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

    written, err := io.Copy(w, response.Body)
    if err != nil {
        log.Printf("Error copying response: %v", err)
        return
    }
    log.Printf("Successfully proxied %d bytes", written)
}

func copyHeaders(dst, src http.Header) {
    for key, values := range src {
        for _, value := range values {
            dst.Add(key, value)
        }
    }
}

func main() {
    proxy := NewProxyServer()

    // Create TLS configuration
    tlsConfig := &tls.Config{
        MinVersion:               tls.VersionTLS12,
        PreferServerCipherSuites: true,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
        },
    }

    server := &http.Server{
        Addr:              ":8443",
        Handler:           http.HandlerFunc(proxy.ProxyHandler),
        ReadTimeout:       timeout,
        WriteTimeout:      timeout,
        ReadHeaderTimeout: readHeaderTimeout,
        MaxHeaderBytes:    1 << 20, // 1 MB
        TLSConfig:        tlsConfig,
    }

    log.Printf("Starting TLS proxy server on port %s", server.Addr)
    if err := server.ListenAndServeTLS("certs/cert.pem", "certs/key.pem"); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}
