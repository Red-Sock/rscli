# RedSock CLI
## rscli - is a simple cli tool for projects and environment handling

### Installation
``` shell
 go install github.com/Red-Sock/rscli@latest
 ```


### Features: 
  - create and manage a Golang projects
  - setup environment for existing projects

### TODO
 - manage cloud infrastructure

### Configuration
Can be used in order to override default settings. All fields are optional to specify

#### Via conf file
You can create rscli.conf file and specify options you'd like to override 
- Either put config next to rscli binary  
```
  $GOPATH/bin/rscli.yaml
```
- or pass config via flag 
```
  rscli [COMMAND] [ARGUMENTS] --rscli-cfg ./rscli.yaml
```

#### Via environment variables
Alternatively (or additionally) you can specify environment variable(s)
that start with "RSCLI_" and followed by field name.

##### Configuration structure (all fields are optional)
- **default_project_git_path** - URL path to git system where this package will be available to fetch
```shell
  export RSCLI_DEFAULT_PROJECT_GIT_PATH=github.com/RedSock
```
- env - object. defines environment handling
  - **path_to_main** - path to main file. Used for project scan. Entrypoint when starts
  ```shell
    export RSCLI_PATH_TO_MAIN=cmd/proj_name/main.go
  ```
  - **path_to_config** - path to project config file. Used for project scan 
  ```shell
    export RSCLI_PATH_TO_CONFIG=config/dev.yaml
  ```

#### Config example can be found in internal/config/rscli.yaml
```yaml
env:
  path_to_main: 'cmd/proj_name/main.go'
  path_to_config: 'config/dev.yaml'

default_project_git_path: github.com/Red-Sock
```
#### SPECIAL VARS
Reading examples you might notice that there is a word **proj_name** used to define a project name. 
This is not just an example. In order to make tool more useful some sequences are predefined for internal use.
- proj_name - can be used in order to put actual name of project somewere
```text
example: 
    When you set variable 
    env.path_to_main to "cmd/proj_name/main.go" 
    and create a project named Rscli via 
    
    rscli create project
    
    it creates project with main file at
    "./github.com/RedSock/rscli/cmd/rscli/main.go"
    
    instead of "./github.com/RedSock/rscli/cmd/proj_name/main.go"
```
### THE LAST, NOT THE LEAST on project creation
- ALL project created via rscli tool require some url at the begging.
- Bear in mind that project is being created withing the whole path, meaning github.com/RedSock/ folders being created


### Environment setup - in development
Something quite interesting...
#### Problem to solve
You are a new guy on a really legacy-full old project with a lot of dependencies on different projects or a data resources.

In order to make your life easier tool scans configuration files 
and projects in current folder and creates environment based on 
what it found
#### Example 
For projects in folder "projects-folder-name" structured as:

```
- .../projects-folder-name/...
  - .../web-application-backend
  - .../web-application-frontend
  - .../web-application-auth
  - .../some-useful-microservice
  - .../another-great-frontend
```
Directory "environment" will be created. 
Environment contains some basic files to setup and connect services:
- docker-compose.example.yaml - Template to be used when updating or creating new envs
- .env.example - environment example file that will be placed next in each project and used during compose
- Makefile / (TODO: PowerShell) - additional scripts to build dev environment. Can be used to define binary build process. 
Migrations and other useful stuff
- Folders with scanned projects:
  - .env - environment variables to configure project (simple merge with .env.example file from above )
  - docker-compose.yaml - compose file based on project dependencies and config scans 
### TODO
  1. [ ]   Local service mesh.
     By specifying in config required external services 
     they will be found and deployed as dependency 
     (like pg or redis but as a service with api)
  2. [ ] Docker networking
  3. [ ] Services mesh with different root path (now service mesh planned to be on a folder level)
  4. [ ] Service mesh with external projects (will be downloaded and launched as standalone applications) 
### Restrictions
- It is only available for projects from the same folder to be meshed
