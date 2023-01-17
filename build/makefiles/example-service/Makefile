APP_NAME ?= example-service
DOCKER_REGISTRY ?= docker.io
DOCKER_REPO ?= example-repo
VERSION ?= $(shell git rev-parse  HEAD | cut -c 1-12)
PWD := $(dir $(abspath $(firstword $(MAKEFILE_LIST))))

IMAGE_TAG?=development


all: release

build: clean
	$(Q)echo 'build example-service $(VERSION) to $(DOCKER_REPO)'
	$(Q)docker build -f ./build/dockerfiles/example-service/Dockerfile --target=runner -t $(APP_NAME) .

test:
	$(Q)echo 'test example-service'
	$(Q)docker build -f ./build/dockerfiles/example-service/Dockerfile --target=tester .

tag:
	$(Q)echo 'tagging example-service $(VERSION) to $(DOCKER_REPO)'
	docker tag $(APP_NAME) $(DOCKER_REGISTRY)/$(DOCKER_REPO)/$(APP_NAME):$(VERSION)

publish: tag
	$(Q)echo 'publish example-service $(VERSION) to $(DOCKER_REPO)'
	docker push $(DOCKER_REGISTRY)/$(DOCKER_REPO)/$(APP_NAME):$(VERSION)

release: ../.. test publish

clean:
	$(Q)echo 'cleaning'

