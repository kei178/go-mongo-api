# go-mongo-api
REST API example with GO & MongoDB

## Requirements

- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/)
- [docker-compose](https://docs.docker.com/compose/)

## Running

```
docker-compose up
```

You can access the API server at `http://localhost:8080`.

If you change something, execute the following command, instead:

```
docker-compose up --build
```

## Shutdown

```
docker-compose down
```

## APIs

The entity Book has the following fields:

- ID (primitive.ObjectID)
- title (string)
- author (string)

Books APIs are as follows:

|METHOD|ENDPOINT|REQUEST HEADERS|REQUEST PAYLOAD|RESPONSE PAYLOAD|
|------|---|---------------|---------------|----------------|
|GET|/books| | |Book[]|
|POST|/books|Content-Type: "application/json"|Book|Create result|
|GET|/books/:id | | |Book|
|PUT|/books/:id |Content-Type: "application/json"|Book|Update result|
|DELETE|/books/:id | | |Delete result| 