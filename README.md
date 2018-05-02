# Fresh8

Fresh8 gaming tech test

## Projects

The repository contains 4 projects

- **apiutil** - library that contains api client utilities
- **feed** - library that contains the feed api client
- **importer** - the importer application
- **store** - library that contains the store api client

## Build & Run

Go to the importer project directory and run:

```
go build && STORE_ADDR=localhost:8001 ./importer
```

Where `localhost:8001` should be replaced with the event store address.

## Configuration

The importer project contains a config.yml file which contains the following settings and defaults:

| Setting                         | Default     |
|---------------------------------|-------------|
| `services.client.timeout`       | `10`        |
| `services.client.feed.host`     | `localhost` |
| `services.client.feed.port`     | `8000`      |
| `jobs.event.steps.limit`        | `10`        |