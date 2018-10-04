IMAGE := akoeb/shoppinglist
NAME := shoppinglist
VERSION := latest
HTTP_USER := liste
HTTP_PW := machen

.PHONY: all run bash rm

all: run

run:
	docker run -d --name $(NAME) -v $(PWD)/shoppinglist.db:/data/shoppinglist.db -e HTTP_USER="$(HTTP_USER)" -e HTTP_PASSWORD="$(HTTP_PW)" -p 8080 --network wordpress $(IMAGE):$(VERSION)

bash:
	docker run -d --name $(NAME) -ti $(IMAGE):$(VERSION) /bin/bash

rm:
	-docker stop $(NAME)
	docker rm $(NAME)
	
