# shellsnake

Play snake in your shell!

<img src="https://user-images.githubusercontent.com/116432/192135742-7fad5354-da71-426e-aa97-3191593e9d89.gif" width="282" alt="A GIF showing a game of shellsnake in progress" />

# To play

## Run from source

Clone this repo, then:

```sh
go run .
```

## Install the binary

```sh
go install github.com/simonwhitaker/shellsnake@latest
```

Then (assuming `$(go env GOPATH)/bin` is in your path):

```sh
shellsnake
```

# Running with Docker

```sh
docker build -t shellsnake .
docker run --rm -it shellsnake:latest
```
