binary:
	# build go binary
	go build -o glockchain main.go

container:
	# build image
	docker build --label glockchain -t glockchain-contained .

clean:
	# delete glockchain binary and db
	-rm glockchain*

	# delete vendor directory
	-rm -rf vendor