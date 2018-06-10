# Builder Container
FROM golang:1.10 as builder
WORKDIR $GOPATH/src/github.com/goosetacob/glockchain/

# Install dep
RUN go get -u github.com/golang/dep/...

# Copy code from host and compile
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./
RUN go build -o /bin/glockchain main.go

# Final Output Container
FROM golang:alpine

# Copy binary to builder container to and final container
COPY --from=builder /bin/glockchain /bin/glockchain

# Run
ENTRYPOINT ["/bin/server"]