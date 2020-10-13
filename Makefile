.PHONY: build release clean all 

all: clean build

build: binary/memdisk-cloudwatch

binary/memdisk-cloudwatch:
	cd docker-build && docker build -t memdisk-cloudwatch-build .
	docker run --rm \
				-v $$(pwd)/binary:/go/bin \
				-v $$(pwd)/src:/go/src/memdisk-cloudwatch \
				memdisk-cloudwatch-build \
				bash -c "cd /go/src/memdisk-cloudwatch && go install -v"

docker_debug: 
	docker run --rm -it \
           -v $$(pwd)/binary:/go/bin \
           -v $$(pwd)/src:/go/src/memdisk-cloudwatch \
           memdisk-cloudwatch-build \
           bash

release: all 
	cd binary && cp memdisk-cloudwatch memdisk-cloudwatch-x86_64 && gzip -9 memdisk-cloudwatch-x86_64

clean: 
	rm -Rf ./binary 

docker-clean:
	docker rmi memdisk-cloudwatch-build

RPM_SOURCEDIR=$(shell rpm --eval "%{_sourcedir}")

rpm-get-sources: build memdisk-cloudwatch.spec memdisk-cloudwatch.service
	rpmdev-setuptree
	cp binary/memdisk-cloudwatch $(RPM_SOURCEDIR)
	cp memdisk-cloudwatch.service $(RPM_SOURCEDIR)

rpm: memdisk-cloudwatch.spec rpm-get-sources
	rpmbuild -ba $(RPMBUILD_FLAGS) $<
