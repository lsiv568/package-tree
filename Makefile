docker_image_tag=package-tree

all: docker-image code

docker-image:
	docker build  -t $(docker_image_tag) .

code:
	docker run -v=`pwd`:`pwd` -w=`pwd` -e GOOS=linux $(docker_image_tag) ./build.sh
