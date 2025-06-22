# Go Rate Limiter: Fixed, Sliding, and Token Bucket

This project demonstrates the implementation of **three popular rate limiting algorithms** in Go that together form a high-performance, multi-user rate limiting system, similar to those used by large-scale APIs such as Twitter or GitHub. The system satisfies the following requirements:

* Supports multiple types of rate limiting:

    * Fixed Window
    * Sliding Window
    * Token Bucket
* Handles concurrent users efficiently
* Designed for performance, with extensibility to use external resources like databases or queues

It also includes a basic HTTP server that applies the rate limiter as middleware, and a client that simulates multiple users sending requests.

---

## Rate Limiting Algorithms

### 1. Fixed Window

* Divides time into fixed intervals (e.g., 1 second).
* Requests are counted per user per interval.
* Simple and fast, but susceptible to bursts at window edges.

### 2. Sliding Window

* Tracks timestamps of recent requests in a sliding time window.
* Fairer distribution and smoother traffic handling.

### 3. Token Bucket

* Each user has a bucket of tokens that refill at a defined rate.
* A request is allowed if a token is available.
* Supports burst traffic and smooths request rate over time.

---

## Running Unit Tests

```bash
go test . -v
```

Each algorithm has corresponding test files:

* `fixed_window_test.go`
* `sliding_window_test.go`
* `token_bucket_test.go`

---

## Example: Running HTTP Server and Client

### Start the HTTP Server

```bash
go run cmd/server/main.go --rate-limiter=token --limit=5
```

| Flag             | Description                                |
| ---------------- | ------------------------------------------ |
| `--rate-limiter` | Choose `fixed`, `sliding`, or `token`      |
| `--limit`        | Max number of requests per second per user |

* The server listens at: `http://localhost:8080`
* Each user is identified via the `X-User-ID` header.

### Simulate Users with Client

```bash
go run cmd/client/main.go --users=10 --requests=6
```

| Flag         | Description                                     |
| ------------ | ----------------------------------------------- |
| `--users`    | Number of simulated concurrent users            |
| `--requests` | Number of requests each user sends              |


Sample output:

```
[User 1, request no 6] (429): Too Many Requests
```

---

