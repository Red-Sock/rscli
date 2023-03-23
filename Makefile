PATTERN_COMPILER_NAME=pattern_compiler

.PHONY: compile-pattern
compile-pattern:
	@echo Compiling project pattern...
	go run support/compiler/main.go

	#rm -rf $(PATTERN_COMPILER_NAME)
	#go build -o $(PATTERN_COMPILER_NAME) support/compiler/main.go
	#./$(PATTERN_COMPILER_NAME)
	#rm -rf $(PATTERN_COMPILER_NAME)
	@echo Project pattern is succesfully compiled!
