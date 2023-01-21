.PHONY: compile-pattern
compile-pattern:
	go run support/compiler/main.go


.PHONY: .compile-plugins
.compile-plugins: compile-pattern