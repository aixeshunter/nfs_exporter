GO ?= go
FIRST_GOPATH := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))
pkgs         = $(shell $(GO) list ./... | grep -v /vendor/)
PREFIX ?= _outputs

DOCKERFILE	  ?= Dockerfile
DOCKER_REPO       ?= aixeshunter
DOCKER_IMAGE_NAME ?= nfs-exporter
DOCKER_IMAGE_TAG  ?= v1.0

PROMU := $(FIRST_GOPATH)/bin/promu

.PHONY: promu
promu:
	GOOS= GOARCH= $(GO) get -v -u github.com/prometheus/promu

.PHONY: build
build: promu
	@echo ">> building binaries"
	$(PROMU) build -v

.PHONY: docker
docker: 
	@echo ">> building docker image from $(DOCKERFILE)"
	@docker build -t "$(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .

gotest: vet format
	@echo ">> running tests"
	@$(GO) test -short $(pkgs)

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)