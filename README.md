# Glockchain
Implementing a simple blockchain with Golang and docker.

Following along with https://jeiwan.cc/posts/building-blockchain-in-go-part-1/ for reference.

# Instructions
1. Install Docker: https://docs.docker.com/install/
2. `$ make run`

# Example (WIP)
```
gustavo at huitzilopochtli in ~/Projects/go/src/github.com/goosetaco/glockchain
$ make run
# build image
docker build -t glockchain-contained .
Sending build context to Docker daemon  11.26kB
Step 1/4 : FROM golang:alpine
 ---> 44ccce322b34
Step 2/4 : WORKDIR $GOPATH/src/github.com/goosetaco/glockchain/
 ---> Using cache
 ---> 4ee3c90c1615
Step 3/4 : COPY ./ $GOPATH/src/github.com/goosetaco/glockchain/
 ---> 024be005d1a7
Step 4/4 : CMD go run main.go
 ---> Running in 77b1b0a0b524
Removing intermediate container 77b1b0a0b524
 ---> c4e6748d089b
Successfully built c4e6748d089b
Successfully tagged glockchain-contained:latest
# run image
docker run glockchain-contained
Mining the block containing "Genesis Block"
0000002a82cc1e17e73782add3553465563453d25a02d249c06988d7da8f1185

Mining the block containing "Send 23 BTC to Brenda"
00000043b99d0b5017acd35e1500da10208ecfc5234e6e58168c53d22ecd9203

Mining the block containing "Send 19 BTC to Pedro"
000000788f83a3faaf0f504d168e182ec5b4d3e7d4f914301565068cb10d21dc

Mining the block containing "Send 11 more BTC to Julio"
000000fcd8f87a137fd5763536d3b39c9d271e4f4668813e8ac67385fa1bc52b

Prev. hash:
Data: Genesis Block
Hash: 0000002a82cc1e17e73782add3553465563453d25a02d249c06988d7da8f1185
PoW: true

Prev. hash: 0000002a82cc1e17e73782add3553465563453d25a02d249c06988d7da8f1185
Data: Send 23 BTC to Brenda
Hash: 00000043b99d0b5017acd35e1500da10208ecfc5234e6e58168c53d22ecd9203
PoW: true

Prev. hash: 00000043b99d0b5017acd35e1500da10208ecfc5234e6e58168c53d22ecd9203
Data: Send 19 BTC to Pedro
Hash: 000000788f83a3faaf0f504d168e182ec5b4d3e7d4f914301565068cb10d21dc
PoW: true

Prev. hash: 000000788f83a3faaf0f504d168e182ec5b4d3e7d4f914301565068cb10d21dc
Data: Send 11 more BTC to Julio
Hash: 000000fcd8f87a137fd5763536d3b39c9d271e4f4668813e8ac67385fa1bc52b
PoW: true
```