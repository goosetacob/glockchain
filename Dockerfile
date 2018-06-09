FROM golang:alpine

WORKDIR $GOPATH/src/github.com/goosetacob/glockchain/

COPY ./ $GOPATH/src/github.com/goosetacob/glockchain/main

CMD ./main