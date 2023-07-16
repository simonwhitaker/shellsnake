# shellsnake

Play snake in your shell!

<img src="https://user-images.githubusercontent.com/116432/192135742-7fad5354-da71-426e-aa97-3191593e9d89.gif" width="282" alt="A GIF showing a game of shellsnake in progress" />

# To play

## Install with Homebrew

```sh
brew install simonwhitaker/tap/shellsnake
```

## Install the binary using the Go toolchain

If you have the Go toolchain installed, you can:

```sh
go install github.com/simonwhitaker/shellsnake@latest
```

## Download a release

Check the [releases page](https://github.com/simonwhitaker/shellsnake/releases). Download, open the gz/zip file, and run `shellsnake`.

## Run from source

Clone this repo, then:

```sh
go run .
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
