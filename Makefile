run:
	# build image
	docker build -t glockchain-contained .

	# run image
	docker run glockchain-contained