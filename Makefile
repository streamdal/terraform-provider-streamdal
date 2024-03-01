HOSTNAME=streamdal.com
NAMESPACE=tf
NAME=streamdal
BINARY=terraform-provider-${NAME}
VERSION=0.1.0
OS_ARCH=darwin_arm64

default: install

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	rm -Rf ./examples/.terraform  ./examples/.terraform.lock.hcl ./examples/terraform.tfstate ./examples/create_pipeline/terraform.tfstate ./examples/create_pipeline/.terraform ./examples/create_pipeline/.terraform.lock.hcl