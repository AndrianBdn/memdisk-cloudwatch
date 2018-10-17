.PHONY: build release clean all 

all: clean build

build:
	cd docker-build && docker build -t memdisk-cloudwatch-build .
	docker run --rm \
				-v $$(pwd)/binary:/go/bin \
				-v $$(pwd)/src:/go/src/memdisk-cloudwatch \
				awsinstancedata-build \
				bash -c "cd /go/src/memdisk-cloudwatch && go install -v"

docker_debug: 
	docker run --rm -it \
           -v $$(pwd)/binary:/go/bin \
           -v $$(pwd)/src:/go/src/memdisk-cloudwatch \
           awsinstancedata-build \
           bash

release: all 
	cd binary && cp memdisk-cloudwatch memdisk-cloudwatch-x86_64 && gzip -9 memdisk-cloudwatch-x86_64

clean: 
	rm -Rf ./binary 

docker-clean:
	docker rmi memdisk-cloudwatch-build

fetch-tarball: memdisk-cloudwatch.spec
	rpmdev-setuptree
	spectool --get-files --sourcedir $<

rpm: memdisk-cloudwatch.spec fetch-tarball
	rpmbuild -ba $(RPMBUILD_FLAGS) $<
