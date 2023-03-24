PATTERN_COMPILER_NAME=pattern_compiler

.PHONY: compile-pattern
compile-pattern:
	@echo Compiling project pattern...
	go run support/compiler/main.go
	@echo Project pattern is succesfully compiled!
