# Template GO Service

This service is a demonstration of using Golang in a monorepo structure.
The service includes a script that can be run to create new services, which will automatically generate :
- Dockerfile
- Makefile
- Github Actions
- Go entrypoint under the "cmd" directory. 
- kubernetes basic deployment & service components

To create a new service, run the following command:
```bash
go run scripts/service-generator/main.go --name=new-service --is-add=true
```

To delete a service, run the following command:
```bash
go run scripts/service-generator/main.go --name=new-service --is-add=false
```
