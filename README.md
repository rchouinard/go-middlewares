# Go Middlewares Project

This project contains a collection of middlewares for Go's `net/http` servers. The middlewares include a logger and a request ID generator.

## Features

- **Logger**: Logs requests and responses in text or JSON format.
- **Request ID Generator**: Adds a unique ID to each request, which can be retrieved from the `X-Request-Id` response header or from the request context.

## Usage

The middlewares can be used in your Go project by importing the `middlewares` package and adding the desired middlewares to your HTTP handler chain.

### Logger

The logger middleware can be added to your HTTP handler chain using the `Logger` or `JSONLogger` functions. These functions take a `http.Handler` as an argument and return a new `http.Handler` that logs requests and responses using the specified format.

```go
import "github.com/rchouinard/go-middlewares"

http.Handle("/", middlewares.Logger(myHandler))
```

### Request ID Generator

The request ID generator middleware can be added to your HTTP handler chain using the `RequestID` function. This function takes a `http.Handler` as an argument and returns a new `http.Handler` that adds a unique ID to each request.

```go
import "github.com/rchouinard/go-middlewares"

http.Handle("/", middlewares.RequestID(myHandler))
```

## Dependencies

This project depends on the following external packages:

- [github.com/oklog/ulid/v2](https://github.com/oklog/ulid): A universally unique lexicographically sortable identifier.
- [github.com/google/uuid](https://github.com/google/uuid): A Go package for generating UUIDs.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
