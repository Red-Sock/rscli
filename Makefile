.PHONY: compile-pattern
compile-pattern:
	go run support/compiler/main.go

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