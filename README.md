# gao
API Development in Go Using Goa

Generate the boilerplate code. The following command takes the design package import path as an argument. It also accepts the path to the output directory as an optional flag:

`goa gen clients/design`

generate a basic implementation of the service along with buildable server files that spin up goroutines to start a HTTP server and client files that can make requests to that server:

`goa example clients/design`

create server and client binaries:
`go build ./cmd/clients`
`go build ./cmd/clients-cli`

run the server: `./clients`