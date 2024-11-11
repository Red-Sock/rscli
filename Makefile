gen-test-project-with-deps: .compile-pattern .gen-test-project-with-deps

.compile-pattern:
	@echo Compiling project pattern...
	go run support/compiler/main.go
	@echo Project pattern is succesfully compiled!

.deps:
	go install github.com/gojuno/minimock/v3/cmd/minimock@latest

mock:
	minimock -i github.com/Red-Sock/rscli/internal/stdio.* -o tests/mocks -g -s "_mock.go"

.gen-test-project-with-deps:
	go build -o rscli-dev
	rm -rf test/testproj
	mkdir -p test
	mkdir -p test/testproj

	cd test &&\
    ./../rscli-dev project init Testproj && \
    cd testproj && \
    ./../../rscli-dev project add grpc

dev-build:
	go build -o $(GOBIN)/rscli-dev .