IMAGE_NAME = todo-app
CONTAINER_NAME = todo-container
PORT = 7540

build:
	docker build -t $(IMAGE_NAME) .

run:
	docker run -d -p $(PORT):$(PORT) --name $(CONTAINER_NAME) $(IMAGE_NAME)

stop:
	docker stop $(CONTAINER_NAME)

rm:
	docker rm $(CONTAINER_NAME)

restart: stop rm build run

exec:
	docker exec -it $(CONTAINER_NAME) /bin/sh

rmi:
	docker rmi $(IMAGE_NAME)