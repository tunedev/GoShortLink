# GoShortLink - URL Shortener

GoShortLink is a simple yet powerful URL shortener service written in Go. The service provides an easy way to shorten long URLs for easier sharing and management. The application is dockerized for smooth development and deployment experience.

## Features

- **Efficient URL shortening:** Shorten any lengthy URLs in a matter of milliseconds.
- **Fast redirection:** Redirect from the shortened URL to the original URL quickly.
- **Persistence:** MongoDB is used as a datastore for long-term URL persistence.

## Getting Started

### prerequisites

- Docker
- Docker Compose
- Go (for locl development withoug Docker)

### Running with Docker Compose

```sh
docker-compose up
```

This will start the GoShortLink service on port 9900 and a MongoDB service on port 27017.

### Building and Runnin Locally

To build the application, run:

```sh
go build -o main ./cmd/server
```

And then to start the server, run

```sh
./main
```

Make sure to set the required environment variables.

### Usage

#### Shorten a URL

Send `POST` request to `http://localhost:9900/shorten` with a JSON body of the following structure:

````JSON
{
  "long_url": "https://example.com/some_long_url/with/long_query_params?short=false"
}```

The service will return a JSON response with the shortened URL:

```JSON
{
  "long_url": "http://localhost:9900/abc123"
}```

#### Redirect to Original URL

Simply navigate to the shortened URL (for example, `http://localhost:9900/abc123`) in a web browser, and you will be redirected to the original URL

### Configuration

The application can be configured using the following environment variables:

* `MONGO_URI`: The MongoDB connection string. Default is `mongodb://localhost:27017`.
* `DB_NAME`: The MongoDB database name. Default is 'goshortlink'
* `COLLECTION_NAME`: The MongoDB collection name. Default is 'urls'
* `PORT`: The port on which the server runs. Default is '9900'

### Contribution

Contributions are always welcome! Please feel free to submit a Pull Request or Open an issue.

License

This project is MIT licensed.
````
