# MyTool

A golang challenge: CLI Swiss Army Knife.

## Usage

`mytool` runs in a Docker Image. You can build your own image today, and soon there will be a public image you can run.

### Requirements

* Docker, to run `mytool`.
* make, to build your own Docker image of `mytool`.
* AWS credentials, for `ecr`.

### Run PreBuilt Docker

This would be easier for users if they could use an existing image.
TODO: Post image to public repository using some gitops automation, Github Actions for example, so thats users don't have to build their own.

### Build and Run

To build the docker image:
`make build-docker`

To run the docker image pass commands and arguments:
` docker run --rm mytool $COMMAND $ARGUMENTS`

For example to run the `ecr` ECR tool, run:
`docker run --rm -v "$HOME/.aws:/root/.aws:rw" mytool ecr --repo repo --image svcImg`

### ECR Digest Finder

Subcommand: `ecr`

Requires AWS credentials mounted via `.aws` in user's home folder. These must be mounted to run in docker.

```
docker run --rm -v "$HOME/.aws:/root/.aws:rw" mytool ecr $IMAGE $TAG 
```
Example. For image `svcImg` and tag `latest`, run:
`docker run --rm -v "$HOME/.aws:/root/.aws:rw" mytool ecr svcImg latest`

### Net Cat

The `nc` subcommand has a server and client mode.
Start the server with `--listen true` in one terminal window.
`docker run --rm -p 8021:8021 mytool nc --listen true --host 192.168.0.2`
Start the client with `--listen false`, the default, in another terminal window.
`docker run --rm mytool nc --host 192.168.0.2`

NOTE: You must change the `docker run` `-p` if using the non-default port.
For example to listen on port `8022` instead:
`docker run --rm -p 8022:8022 mytool nc --listen true --port 8022`

To run without docker on the default ports, you can also use `go run main.go nc --listen true` for the server,
and `go run main.go nc` for the client.

## Developing mytool

Working on `mytool` requires:
* go 1.17+

To run using go:
```
 go run main.go $COMMAND $ARGS
```
ie
```
go run main.go ecr --image woe-sim --repo woe-sim
```
New commands can be added with the `cobra` templating.

