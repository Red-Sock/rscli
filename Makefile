build:
	ifeq ($(OS),Windows_NT)
		echo "windows"# go build -o ./build/rscli.exe
	else
		go build -o ./build/rscli.exe
	endif

.run:
run:
	go run main.go
