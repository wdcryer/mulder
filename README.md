# The Jenkins X Files - Mulder

`mulder` is the backend side of [The Jenkins X Files](https://the-jenkins-x-files.github.io/) - the [Jenkins X](https://jenkins-x.io/) workshop. You can also see [Scully](https://github.com/the-jenkins-x-files/scully), the frontend side.

It's a Go application that provides a (very) basic HTTP API, with 1 main endpoint:

- `GET /quote/random` which returns a random quote from FBI's most unwanted, in JSON:

    ```
    {
        "quote": "I have a theory. Do you want to hear it?"
    }
    ```

- `GET /healthz` checks the health of the application, and the connection to Redis. It returns either a `200` or `500` status code.

**Dependencies**:

- [Redis](https://redis.io/) - to store the quotes

**Building**:

- `go build`

**Running**:

Either:
- build the binary with `go build`, and run it
- or run it directly with `go run .`

**Flags**:

- `-listen-addr` (string): host:port on which to listen. Default: `:8080`
- `-redis-addr` (string): redis host:port to connect to. Default: `:6379`
- `-redis-connect-timeout` (duration): timeout for connecting to redis. Default: `1m0s`

**Unit Tests**:

- `go test -v .`

**Integration Tests**:

- `go test -v ./tests -addr HOST:PORT`
    - don't forget to replace the `HOST:PORT` argument with the hostname and port of a running mulder instance you want to test - the integration tests won't start it for you.