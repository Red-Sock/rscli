# RS CLI
CLI tool to handle golang projects

Allows you to create and manage a Golang project

## Installing
go install github.com/Red-Sock/rscli@latest

## Already supported features

- Project build
- Separated configuration build
- Environment setup

### Project and configuration build
can be used to create and configure an application - basically a boilerplate to a fast start

### Environment setup - in development
Something quite interesting...

For projects in folder "projects-folder-name" structured as:
- .../projects-folder-name/...
  - .../web-application-backend
  - .../web-application-frontend
  - .../web-application-auth
  - .../some-useful-microservice
  - .../another-great-frontend

creates a directory "environment" that contains some basic files to setup and connect all of the above:
- docker-compose.example.yaml - template to be used when updating or creating new envs
- .env.example - environment example file that will be placed next to each project
- Makefile / PowerShell (TODO) - additional scripts to build dev environment
- folders with projects:
  - .env - with environment variables to configure project
  - docker-compose.yaml - with all needed resources



#### Terminal UI is based on [rscli-uikit](https://github.com/Red-Sock/rscli-uikit)