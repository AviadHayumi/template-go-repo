clean:
	make  -f build/makefiles/*/Makefile clean

build: clean
	make  -f build/makefiles/*/Makefile build

test:
	make  -f build/makefiles/*/Makefile test

tag:
	make  -f build/makefiles/*/Makefile tag

publish:
	make  -f build/makefiles/*/Makefile publish

run-compose:
	cd build/docker-compose && docker-compose up

run-kubernetes:
	kubectl apply -f build/kubernetes

run-kubernetes-minikube:
	eval $$(minikube docker-env)
	make build tag VERSION=latest
	kubectl apply -f build/kubernetes