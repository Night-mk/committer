prepare:
	@sudo rm -rf /tmp/badger
	@go build
	@#cd examples/client && go build
	@cd controller && go build
	@mkdir /tmp/badger
	@mkdir /tmp/badger/coordinator
	@mkdir /tmp/badger/follower
	@mkdir /tmp/badger/follower/3001
	@mkdir /tmp/badger/follower/3002
	@mkdir /tmp/badger/follower/3003
	@mkdir /tmp/badger/follower/3004
	@mkdir /tmp/badger/follower/3005
	@mkdir /tmp/badger/follower/3006


run-example-coordinator:
	@./committer -role=coordinator -nodeaddr=localhost:3000 -followers=localhost:3001 -committype=three-phase -timeout=1000 -dbpath=/tmp/badger/coordinator -whitelist=127.0.0.1

run-example-follower:
	@./committer -role=follower -nodeaddr=localhost:3001 -committype=three-phase -timeout=1000 -dbpath=./badger -whitelist=127.0.0.1

run-example-client:
	@examples/client/client

unit-tests:
	@cd server && go test

functional-tests:
	@go test