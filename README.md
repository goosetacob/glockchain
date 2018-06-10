# Glockchain
Implementing a simple blockchain with Golang and docker.

Following along with https://jeiwan.cc/posts/building-blockchain-in-go-part-1/ for reference.

### How to use
*still wip so may change may change before i can update this*
```
gustavo at huitzilopochtli in ~/Projects/go/src/github.com/goosetacob/glockchain (master)
$ make binary
# build go binary
go build -o glockchain main.go

gustavo at huitzilopochtli in ~/Projects/go/src/github.com/goosetacob/glockchain (master)
$ ./glockchain
Usage of balance:
  -address string
        The address to send genesis block reward to
Usage of create:
  -address string
        The address to send genesis block reward to
Usage of send:
  -amount int
        Amount to send
  -from string
        Source wallet address
  -to string
        Destination wallet address

gustavo at huitzilopochtli in ~/Projects/go/src/github.com/goosetacob/glockchain (master)
$ ./glockchain create -address "gustavo"
INFO[0000] Mining the block containing "[ID: 4f953d1895e67f08f7b2b624dff2de33102d532041dad813efeae789456c779b Inputs: [  -1 GooseCoin `all I want is to just have fun live my life like a son of a gun`] Outputs: [ gustavo : 1000000]]"
INFO[0047] Done!
```