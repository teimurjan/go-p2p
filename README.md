# About

go-p2p - is a peer-to-peer network targeting on the local network in order to have a distributed file system inside the companies.

# Installation

## Docker

The app is dockerized so it can be started using [Docker](https://docs.docker.com/docker-for-mac/install/) and [Docker Compose](https://docs.docker.com/compose/install/).

If you have the tools installed on your machine, create `.env` file by example file `.env.example` and type:

### Development

```sh
docker-compose -f docker-compose.dev.yml up -d
```

### Production

```sh
docker-compose -f docker-compose.prod.yml up -d
```

## Locally

If you need to run the app without docker, you will need [Go](https://golang.org/doc/install), [Dep](https://golang.github.io/dep/) and [Node](https://nodejs.org/en/download/).

After the tools installation:

### Backend setup

```sh
dep ensure
```

to install all Go dependencies.

Create `.env` file like in [Docker](##Docker) section and run it:

```sh
go run main.go
```

### GUI setup

```sh
cd gui
npm i
npm run start
```
