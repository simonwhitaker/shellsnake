# Emoji Snake

A simply snake game written using [Bubbletea - a powerful little TUI framework](https://github.com/charmbracelet/bubbletea).

<img src="https://user-images.githubusercontent.com/116432/184492979-b12266d5-6b8b-475a-84be-8e47f9f3ccde.gif" width="410" alt="A GIF showing a game of emojisnake in progress" />

# To play

## Run from source

Clone this repo, then:

```sh
go run .
```

## Install the binary

```sh
go install github.com/simonwhitaker/emojisnake@latest
```

Then (assuming `$(go env GOPATH)/bin` is in your path):

```sh
emojisnake
```

# Running with Docker

```sh
docker build -t emojisnake .
docker run --rm -it emojisnake:latest
```
