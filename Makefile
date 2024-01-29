mockgen:
	sh ./scripts/mockgen.sh
bench:
	go test -bench -benchmem -memprofile=./base.out -benchtime=100000x $(go list ./... | grep -v mock')