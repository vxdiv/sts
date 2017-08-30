all: config dep build
config:
	eval "echo \"`cat config.yaml.dist`\"" > config.yaml
dep:
	glide install
build:
	go build -o sts -v
test:
	go test -test.cover -test.race `glide novendor`
bench:
	go test -test.cover -test.race -test.bench=. -test.benchmem `glide novendor`
run:
	go run main.go
