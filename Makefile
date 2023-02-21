export GO111MODULE=on
export CGO_ENABLED=0
export GOFLAGS=-buildvcs=false

check-wire:
	@echo 'check which wire'
	which wire || (go install github.com/google/wire/cmd/wire@latest)

wire: check-wire
	@echo 'generate wire'
	cd ./cmd/main && wire

