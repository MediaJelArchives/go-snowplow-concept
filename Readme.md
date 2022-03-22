# Description

Proof of concept on the snowplow golang tracker

# Setup using `Golang`

Ensure you have `Go >= 1.18` installed. You may follow these [instructions](https://go.dev/doc/install) to install Go in your environment
Ensure you have `PORT 3000` open

```bash
cd go-snowplow-concept

go mod download

go run main.go

# Visit localhost:3000

```

# Setup using `docker-compose`

Ensure you have `docker` and `docker-compose` installed
Ensure you have `PORT 3000` open

```bash
cd go-snowplow-concept

docker-compose up -d

# Visit localhost:3000
```