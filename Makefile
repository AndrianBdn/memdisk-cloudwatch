.PHONY: build release clean all

all: clean build

build:
	cd docker-build && docker build --platform linux/amd64  -t memdisk-cloudwatch-build-x86_64 .
	cd docker-build && docker build --platform linux/arm64  -t memdisk-cloudwatch-build-arm64 .

	docker run --rm --platform linux/amd64 \
				-v $$(pwd)/binary_x86_64:/go/bin \
				-v $$(pwd)/src:/go/src/memdisk-cloudwatch \
				memdisk-cloudwatch-build-x86_64 \
				bash -c "cd /go/src/memdisk-cloudwatch && go install -v"

	docker run --rm --platform linux/arm64 \
				-v $$(pwd)/binary_arm64:/go/bin \
				-v $$(pwd)/src:/go/src/memdisk-cloudwatch \
				memdisk-cloudwatch-build-arm64 \
				bash -c "cd /go/src/memdisk-cloudwatch && go install -v"

docker_debug:
	docker run --rm -it --platform linux/amd64 \
           -v $$(pwd)/binary:/go/bin \
           -v $$(pwd)/src:/go/src/memdisk-cloudwatch \
           memdisk-cloudwatch-build \
           bash

pack:
	cd ./binary_x86_64 && cp memdisk-cloudwatch memdisk-cloudwatch-x86_64 && gzip -9 memdisk-cloudwatch-x86_64
	cd ./binary_arm64 && cp memdisk-cloudwatch memdisk-cloudwatch-arm64 && gzip -9 memdisk-cloudwatch-arm64

release: all pack

clean:
	rm -Rf ./binary_arm64 ./binary_x86_64

docker-clean:
	docker rmi memdisk-cloudwatch-build


