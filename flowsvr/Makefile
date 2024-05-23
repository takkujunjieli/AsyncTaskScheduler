PLATFORM=$(shell uname -m)
DATETIME=$(shell date "+%Y%m%d%H%M%S")

flowsvr:
	go build -o bin/flowsvr src/main.go
	chmod +x bin/flowsvr
test:
	go test -v ./... -coverprofile=cover.out
clean:
	$(RM) tmp/* $(TARGET) 
