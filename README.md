# gohash
Hashes inputs with SHA-512 and returns the hashed value in base64 encoding after 5 seconds.  The inputs should be sent with key values using the application/x-www-form-urlencoded content-type.

Example curl: `curl --data "password=angryMonkey" http://localhost:8080`

Additionally, you can request the server gracefully shut down by sending the input key `graceful shutdown` with the value `true`.  After the server shut down has 'started' (which occurs after a 5 second wait time before the server responds, as it is for hashing requests), the server allows pending requests to finish but will not accept any further requests.

Example graceful shutdown curl: `curl --data "graceful shutdown=true" http://localhost:8080`

## Install
1. Make sure you have [golang installed](https://golang.org/dl/) and clone this repository to your local Go workspace.
2. Run `go build` in the project root directory to generate the executable and run it on the command line  with `./gohash`.
3. The default port used is `8080` but you can specify a different one with `./gohash -port 9000`, for example.

## Testing
You can run the included tests using `go test` from the project root directory.  Currently there are only a couple of minimal tests included as more focus was given on learning to work with Golang over completeness of testing.
