package main

import (
  "io"
  "log"
  "net/http"
)

func proxyHandler(w http.ResponseWriter, r * http.Request) {
  log.Printf("Received Incoming Request: %s %s", r.Method, r.URL)

  outgoingReq, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
  if err != nil {
    log.Printf("Error creating new request: %s", err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  outgoingReq.Header = r.Header

  client := &http.Client{}

  response, err := client.Do(outgoingReq)

  if err != nil {
    log.Printf("Error sending request: %s", err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  defer response.Body.Close()

  for key, values := range response.Header {
    for _, value := range values {
      w.Header().Add(key, value)
    }
  }

  w.WriteHeader(response.StatusCode)
  io.Copy(w, response.Body)
}

func main() {
  port := ":8080"
  http.HandleFunc("/", proxyHandler)
  log.Printf("Starting server on port %s", port)

  if err := http.ListenAndServe(port, nil); err != nil {
    log.Fatalf("Error starting server: %s", err)
  }
}
