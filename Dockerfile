FROM golang:alpine

WORKDIR $GOPATH/src/github.com/goosetaco/glockchain/

COPY ./ $GOPATH/src/github.com/goosetaco/glockchain/

CMD go run main.go