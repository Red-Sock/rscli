# RS CLI
CLI tool

## Configuration

## Project 

### Project layout

1. `/config`

    Holds all configuration files to application:
   1. local-config.yml - for local testing. Excluded by .gitignore
   2. config.yml - main file. Default config for app  
   3. config.yml.example - example to recreate local config

2. `/cmd`

    Holds folders with main files
    e.g:
   1. `/cmd/cron/main.go` - some cron job app
   2. `/cmd/api/grpc/main.go` - some grpc server
   3. `/cmd/api/web/main.go` - some http server

3. `/transport`
   
    Holds all server related stuff:
   1. grpc server structures
   2. http server structures

4. `/internal`

    Holds all internal stuff such as:
   1. client connections (to another services)
   2. business logic layers
   3. data source (db, cache, file system)

5. `/pkg`

   Holds all logic that can be exported e.g:
   1. `/pkg/api/*grpc_generated_files*`
   2. `/pkg/swagger/*swagger_generated_files*`