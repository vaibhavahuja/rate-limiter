## Rate Limiter Service

#### Tech Doc : [here](https://docs.google.com/document/d/1v14bOs-zXegk_AmMWduAnNVyI_Lus2q_ahYjxFAQe8E/edit#heading=h.10r0locdkasw)
The tech doc contains the HLD/LLD and other implementation details.

#### Running the service

- Start docker containers
  - `docker-compose up`
- Run main file
  - `go run cmd/main.go`
- Proto definitions are in this location 
  - ./proto/rate_limiter.proto

#### Running the tests

- Start docker containers 
  - `docker-compose up`
- Run the tests 
  - `go test ./... -v -count=1 -covermode=atomic -coverpkg=./... -coverprofile=coverage.out
    `
- Display Coverage
  - `go tool cover -html=coverage-functional.out -o coverage.html
    `

#### How to contribute?
- Raise an issue
- Talk to any of the contributors