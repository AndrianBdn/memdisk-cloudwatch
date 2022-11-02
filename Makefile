.PHONY: build release clean all

all: clean build

build:
	cd docker-build && docker build -t memdisk-cloudwatch-build .

	docker run --rm \
				-v $$(pwd)/binary:/go/bin \
				-v $$(pwd)/src:/go/src/memdisk-cloudwatch \
				-e GOOS=linux \
				-e GOARCH=amd64 \
				memdisk-cloudwatch-build \
				bash -c "cd /go/src/memdisk-cloudwatch && go build -o /go/bin/memdisk-cloudwatch-x86_64"

	docker run --rm \
				-v $$(pwd)/binary:/go/bin \
				-v $$(pwd)/src:/go/src/memdisk-cloudwatch \
				-e GOOS=linux \
				-e GOARCH=arm64 \
				memdisk-cloudwatch-build \
				bash -c "cd /go/src/memdisk-cloudwatch && go build -o /go/bin/memdisk-cloudwatch-arm64"

docker_debug:
	docker run --rm -it --platform linux/amd64 \
           -v $$(pwd)/binary:/go/bin \
           -v $$(pwd)/src:/go/src/memdisk-cloudwatch \
           memdisk-cloudwatch-build \
           bash

pack:
	cd ./binary &&  gzip -9 memdisk-cloudwatch-x86_64 && gzip -9 memdisk-cloudwatch-arm64

release: all pack

clean:
	rm -Rf ./binary

docker-clean:
	docker rmi memdisk-cloudwatch-build
