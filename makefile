
TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=perxtech.com
NAMESPACE=tf
NAME=metabase
BINARY=terraform-provider-${NAME}
VERSION=0.3
OS := $(shell uname -s | tr A-Z a-z)
ARCH := amd64

default: install

build:
	go build -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_darwin_amd64
	# GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_freebsd_386
	# GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_freebsd_amd64
	# GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_freebsd_arm
	# GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_linux_amd64
	# GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_linux_arm
	# GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_openbsd_386
	# GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_openbsd_amd64
	# GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_solaris_amd64
	# GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_windows_386
	# GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_windows_amd64

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS}_${ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS}_${ARCH}

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

testcov:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m -coverprofile=coverage.out
	go tool cover -html=coverage.out

init:
	cd examples && terraform init

plan: init
	cd examples && terraform plan

apply:
	cd examples && terraform apply
