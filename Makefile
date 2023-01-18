PROJECT_TMP_FOLDER="temp_project"
PROJECT_PATTERN_SRC_REPO=https://github.com/Red-Sock/project-plugin
PROJECT_COMPILED_FOLDER=plugins/src/project/processor/patterns/pattern_c

.PHONY: compile-pattern
compile-pattern:
	echo 'recreating tmp folder $(PROJECT_TMP_FOLDER)'
	echo off
	rm -rf $(PROJECT_TMP_FOLDER)
	mkdir $(PROJECT_TMP_FOLDER)
	echo 'cloning from $(PROJECT_PATTERN_SRC_REPO) to $(PROJECT_TMP_FOLDER)'
	cd $(PROJECT_TMP_FOLDER) && git clone $(PROJECT_PATTERN_SRC_REPO) . &&  git pull && git switch go
	echo 'compiling project to '
	go run support/project-compiler/main.go $(PROJECT_TMP_FOLDER) $(PROJECT_COMPILED_FOLDER)

.PHONY: compile-project-plugin
compile-project-plugin:
	go build -buildmode=plugin -o plugins/project.so plugins/src/project/main.go
.PHONY: compile-project-plugin-ui
compile-project-plugin-ui:
	go build -buildmode=plugin -o plugins/project-ui.so plugins/src/project/ui/main.go

.PHONY: compile-config-plugin
compile-config-plugin:
	go build -buildmode=plugin -o plugins/config.so plugins/src/config/main.go

.PHONY: compile-config-plugin-ui
compile-config-plugin-ui:
	go build -buildmode=plugin -o plugins/config-ui.so plugins/src/config/ui/main.go


.PHONY: .compile-plugins
.compile-plugins: compile-pattern compile-project-plugin compile-project-plugin-ui compile-config-plugin compile-config-plugin-ui