.PHONY: compile-pattern
compile-pattern:
	@echo Compiling project pattern...
	go run support/compiler/main.go
	@echo Project pattern is succesfully compiled!


.deps:
	go install github.com/gojuno/minimock/v3/cmd/minimock@latest

mock:
	minimock -i github.com/Red-Sock/rscli/internal/stdio.* -o tests/mocks -g -s "_mock.go"

testproj:
	cd test &&\
	rm -rf testproj &&\
    go run ./../main.go project init -n testproj &&\
    cd testproj &&\
    go mod tidy
