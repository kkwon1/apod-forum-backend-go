# APOD Forum Backend

Link to APOD Forum application :ringed_planet: -> https://apod.kkwon.dev/

## Introduction

NASA offers an API for a daily media upload related to astronomy called [Astronomy Picture of the Day (APOD)](https://data.nasa.gov/Space-Science/Astronomy-Picture-of-the-Day-API/ez2w-t8ua).

- [NASA's official APOD Github](https://github.com/nasa/apod-api)
- [NASA's official APOD page](https://apod.nasa.gov/apod/astropix.html)

This Go service is a wrapper around the APOD API, which serves an APOD Forum Frontend,
allowing users to retrieve batches of paginated APODs, search for APODs, and interact with posts by liking, saving and commenting. The forum is heavily inspired by [lobste.rs](https://lobste.rs/)

Check out the [APOD Forum Frontend repository](https://github.com/kkwon1/apod-forum-frontend).

## Running the Service

### Prerequisites

This project requires Go to be installed in your local machine. Please visit the official [go.dev](https://go.dev/doc/install) website to download the version you need.

### Local Dev

#### Build

Go comes with its own build toolchain which is baked into the language itself. It handles compilation, dependency management, and testing.

To build the executable file, run

```
> go build ./cmd/main.go
```

#### Run

Or you can choose to build and run the service all in one step.

To run this service locally, run the following command from the root directory of project. The service should run on port `8080`.

```
> go run ./cmd/main.go
```

#### Environment Varialbes

The project loads all environment variables from a `.env` file, which specifies the NASA API Key, MongoDB endpoint, Auth0 token issuer, etc... Please set the appropriate values for running this service locally.

### Docker

To run this service on Docker, run the following commands from the root directory of project. The service should run on port `8080`.

```
> docker-compose build
> docker-compose up
```

Once the container is running, you can ssh into the container by using the command

```
> docker exec -it <CONTAINER_NAME> sh
```

## Testing

To run all tests in current directory and all subdirectories

```
> go test ./...
```

To run a test in a specific directory

```
> cd <directory>
> go test
```
