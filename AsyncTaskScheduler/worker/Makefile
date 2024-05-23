PLATFORM=$(shell uname -m)
DATETIME=$(shell date "+%Y%m%d%H%M%S")

worker:
	go build -o bin/worker src/main.go
	chmod +x bin/worker
test:
	go test -v ./... -coverprofile=cover.out
clean:
	$(RM) tmp/* $(TARGET) 
