.PHONY: gen
gen:
	protoc --go_out=. --go_opt=module=github.com/rphmauriciodev/gopher-ledger \
		--go-grpc_out=. --go-grpc_opt=module=github.com/rphmauriciodev/gopher-ledger \
		proto/ledger.proto

.PHONY: tidy
tidy:
	go mod tidy