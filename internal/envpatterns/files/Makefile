PROJECT_NAME=proj_name
PROJECT_PATH=abs_proj_path
PATH_TO_MAIN_GO=path_to_main

.PHONY: build-app
build-app:
	cd $(PROJECT_PATH) && GOOS=linux go build -o $(PWD)/$(PROJECT_NAME) $(PROJECT_PATH)/$(PATH_TO_MAIN_GO)

.PHONY: service-up
service-up:	build-app env-up app-up

.PHONY: env-up
env-up:
	docker-compose --env-file .env -f docker-compose.yaml -p $(PROJECT_NAME) up -d

.PHONY: env-down
env-down:
	docker-compose -p $(PROJECT_NAME) down

.PHONY: app-up
app-up:
	docker-compose --env-file .env -f docker-compose.yaml -p $(PROJECT_NAME) up -d proj_name