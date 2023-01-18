package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

/*
this script suppose to generate all service boilerplate
*/
func main() {
	logger := logrus.New()
	serviceName := flag.String("name", "hello", "Name of the service")
	isAdd := flag.Bool("is-add", false, "is adding service false gonna delete service")

	flag.Parse()

	if serviceName == nil || *serviceName == "" {
		logger.Fatal("Service name is required")
	}

	if *isAdd {
		logger.Info("Validating service not exists")
		if isServiceAlreadyExists(*serviceName) {
			logger.Fatalf("Service %s already exists", *serviceName)
		}

		logger.Info("Creating service")
		err := createService(*serviceName)
		if err != nil {
			logger.Fatalf("cannot create service %v", err)
		}
		logger.Info("Created service successfully")
	} else {
		deleteService(*serviceName)
	}

}

func deleteService(serviceName string) {
	os.RemoveAll("cmd/" + serviceName)
	os.Remove(fmt.Sprintf(".github/workflows/%v.yaml", serviceName))
	os.RemoveAll(fmt.Sprintf("cmd/%v", serviceName))
	os.RemoveAll("build/makefiles/" + serviceName)
	os.RemoveAll("build/dockerfiles/" + serviceName)
	os.RemoveAll("deployment/kubernetes/" + serviceName)
}

func isServiceAlreadyExists(serviceName string) bool {
	// check if .github/workflow/serviceName exists
	return fileExists(".github/workflow/"+serviceName) || fileExists("cmd/"+serviceName) || fileExists("deployment/kubernetes/"+serviceName)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func createService(serviceName string) error {
	if err := createGithubWorkflow(serviceName); err != nil {
		return err
	}
	if err := createMakefile(serviceName); err != nil {
		return err
	}
	if err := createDockerfile(serviceName); err != nil {
		return err
	}
	if err := createGoEntrypoint(serviceName); err != nil {
		return err
	}
	if err := createKubernetes(serviceName); err != nil {
		return err
	}

	//AddToCompose(serviceName)
	return nil
}

func createKubernetes(serviceName string) error {
	deploymentTemplate := `apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: example-service
  name: example-service
spec:
  replicas: 1
  selector:
    matchLabels:
      app: example-service
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: example-service
    spec:
      containers:
      - image: example-repo/example-service
        imagePullPolicy: IfNotPresent
        name: example-service
        resources: {}
status: {}
`
	svcTemplate := `apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  labels:
    app: example-service
  name: example-service
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: example-service
status:
  loadBalancer: {}
`
	deploymentContent := strings.Replace(deploymentTemplate, "example-service", serviceName, -1)
	svcContent := strings.Replace(svcTemplate, "example-service", serviceName, -1)

	err := os.MkdirAll(fmt.Sprintf("deployment/kubernetes/%v", serviceName), os.ModePerm)
	if err != nil {
		return err
	}
	// save content to file
	err = os.WriteFile(fmt.Sprintf("deployment/kubernetes/%v/deployment.yaml", serviceName), []byte(deploymentContent), 0666)
	err = os.WriteFile(fmt.Sprintf("deployment/kubernetes/%v/service.yaml", serviceName), []byte(svcContent), 0666)
	return err
}

func createGoEntrypoint(serviceName string) error {
	goEntrypointTemplate := `package main

import "fmt"

func main() {
	fmt.Println("example service")
}
`
	goEntryPointContent := strings.Replace(goEntrypointTemplate, "example service", serviceName, -1)

	err := os.MkdirAll("cmd/"+serviceName, os.ModePerm)
	if err != nil {
		return err
	}
	// save content to file
	err = os.WriteFile(fmt.Sprintf("cmd/%v/main.go", serviceName), []byte(goEntryPointContent), 0666)
	return err
}

func createGithubWorkflow(serviceName string) error {
	templateWorkflow := `# This is a basic workflow to help you get started with Actions
name: CI

# Controls when the workflow will run
on:
  # Triggers the workflow on push or pull request events but only for the "main" branch
  push:
    branches: [ "*" ]
#  pull_request:
#    branches: [ "*" ]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v2

      - name: Set up Docker
        uses: docker-practice/actions-setup-docker@master
        timeout-minutes: 20

      - name: Build
        run: make  -f build/makefiles/example-service/Makefile build

      - name: Test
        run: make  -f build/makefiles/example-service/Makefile test

      - name: Tag
        run: make  -f build/makefiles/example-service/Makefile tag

      - name: Publish
        run: echo "publish" #make  -f build/makefiles/example-service/Makefile publish`
	workflowContent := strings.Replace(templateWorkflow, "example-service", serviceName, -1)

	// save content to file
	err := os.WriteFile(fmt.Sprintf(".github/workflows/%v.yaml", serviceName), []byte(workflowContent), 0666)
	return err
}

func createMakefile(serviceName string) error {
	makefileTemplate := `APP_NAME ?= example-service
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
`
	makefileContent := strings.Replace(makefileTemplate, "example-service", serviceName, -1)

	err := os.MkdirAll(fmt.Sprintf("build/makefiles/%v", serviceName), os.ModePerm)
	if err != nil {
		return err
	}
	// save content to file
	err = os.WriteFile(fmt.Sprintf("build/makefiles/%v/Makefile", serviceName), []byte(makefileContent), 0666)
	return err
}

func createDockerfile(serviceName string) error {
	dockerfileTemplate := `FROM golang:alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o /usr/bin/example-service ./cmd/example-service/main.go

FROM builder as tester
RUN go test ./...

FROM alpine as runner
COPY --chown=0:0 --from=builder /usr/bin/example-service /usr/bin
ENTRYPOINT exec example-service
`
	dockerfileContent := strings.Replace(dockerfileTemplate, "example-service", serviceName, -1)

	err := os.MkdirAll("build/dockerfiles/"+serviceName, os.ModePerm)
	if err != nil {
		return err
	}
	// save content to file
	err = os.WriteFile(fmt.Sprintf("build/dockerfiles/%v/Dockerfile", serviceName), []byte(dockerfileContent), 0666)
	return err
}

// AddToCompose todo
func AddToCompose(serviceName string) {

}
