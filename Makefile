build:
	# build binary
	go build main.go

	# build image
	docker build --label glockchain -t glockchain-contained .

run:
	# run image
	docker run glockchain-contained 

debug:
	docker exec -it $(docker ps --filter "label=glockchain" -q) bash

exec:
	docker exec -it glockchain-contained /bin/bash

clean:
	docker kill $(docker ps --filter "label=glockchain" -q)

murder:
	# kill the container
	docker kill "$(docker ps --filter "label=glockchain" -q)"

	# delete db
	rm glockchain_ledger

	# delete bin
	rm main