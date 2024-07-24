package main

import (
    "context"
    "encoding/base64"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "io"
    "net/http"
    "net/http/httptest"
    "strings"
)

type readCloserWrapper struct {
    io.Reader
}

func (rcw readCloserWrapper) Close() error {
    return nil
}

func handleLambdaRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // Create a new HTTP request
    r, err := http.NewRequest(request.HTTPMethod, request.Path, nil)
    if err != nil {
        return events.APIGatewayProxyResponse{StatusCode: 500}, err
    }

    // Add query parameters
    q := r.URL.Query()
    for k, v := range request.QueryStringParameters {
        q.Add(k, v)
    }
    r.URL.RawQuery = q.Encode()

    // Add headers
    for k, v := range request.Headers {
        r.Header.Add(k, v)
    }

    // Handle body
    if request.Body != "" {
        decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(request.Body))
        r.Body = readCloserWrapper{decoder}
    }

    // Create a response recorder
    w := httptest.NewRecorder()

    // Serve the request
    http.DefaultServeMux.ServeHTTP(w, r)

    // Prepare the response
    resp := events.APIGatewayProxyResponse{
        StatusCode: w.Code,
        Headers:    make(map[string]string),
    }

    for k, v := range w.Header() {
        resp.Headers[k] = v[0]
    }

    resp.Body = w.Body.String()

    return resp, nil
}

func startLambda() {
    lambda.Start(handleLambdaRequest)
}