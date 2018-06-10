build:
	# build image
	docker build --label glockchain -t glockchain-contained .

clean:
	# delete glockchain_ledger
	-rm glockchain.db

	# delete vendor directory
	-rm -rf vendor