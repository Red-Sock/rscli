.PHONY: compile-pattern
compile-pattern:
	@echo Compiling project pattern...
	go run support/compiler/main.go
	@echo Project pattern is succesfully compiled!

.deps:
	go install github.com/gojuno/minimock/v3/cmd/minimock@latest

mock:
	minimock -i github.com/Red-Sock/rscli/internal/stdio.* -o tests/mocks -g -s "_mock.go"

gen-test-project-with-deps:
	go build -o rscli-dev
	cd test &&\
	rm -rf testproj &&\
    ./../rscli-dev project init Testproj && \
    cd testproj && \
    ./../../rscli-dev project add postgres redis grpc rest telegram sqlite




dev-build:
	go build -o $(GOBIN)/rscli-dev .