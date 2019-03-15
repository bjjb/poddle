TAG=poddle # the image name

.PHONY: build run

# build the docker image
build:
	docker build -t $(TAG) .

# run the docker image
run: build
	docker run -it -dP -v$(PWD)/app:/srv/app $(TAG)
