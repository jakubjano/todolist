# Onboarding task - TODO list

## Description
Simple microservice architecture project. 

## User Service
* Use gRPC gateway (https://github.com/grpc-ecosystem/grpc-gateway) - API contracts
* Write a docker file to run application
* API contracts: CRUD with token verification using firebase auth (https://firebase.google.com/docs/admin/setup)
* Storage: NoSQL database (Firestore) (https://firebase.google.com/docs/firestore/quickstart)
* Cache: Storing user role for faster access (https://github.com/muesli/cache2go)
* Testing: create unit and integration test (https://github.com/stretchr/testify)

## Task Service
* Use gRPC gateway (https://github.com/grpc-ecosystem/grpc-gateway) - API contracts
* Use HTTP communication with Echo server (https://echo.labstack.com/) - Internal operations
* Write a docker file to run application
* API contracts: CRUD (Authorization) + some basic filter (e.g. timestamp, type etc.)
* Storage: NoSQL database (Firestore) (https://firebase.google.com/docs/firestore/quickstart)
* Testing: create unit and integration test (https://github.com/stretchr/testify)

## Configurability
* Parametrize variables via configuration using Viper (https://github.com/spf13/viper)

## Logging
* Logs have been an essential part of troubleshooting application
* Add logs to critical parts of your application (https://github.com/uber-go/zap)
* Consider adding some information logs.
* log with custom fields 

## Deployment
* Write terraform to manage cloud services (https://www.terraform.io/)
* Each service will use Cloud Run
* Create Cloud storage to store exported Task for specific user's.
* Deploy services to GCP
