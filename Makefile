IMAGE := akoeb/shoppinglist
NAME := shoppinglist
VERSION := latest
HTTP_USER := liste
HTTP_PW := machen

.PHONY: all run bash rm

all: runlocal

compile:
	go build .

runlocal: compile
	foundation watch &
	./shoppinglist -db ./shoppinglist.db

run:
	sudo docker run -d --name $(NAME) -v $(PWD)/shoppinglist.db:/data/shoppinglist.db -e HTTP_USER="$(HTTP_USER)" -e HTTP_PASSWORD="$(HTTP_PW)" -p 8080 --network wordpress $(IMAGE):$(VERSION)

bash:
	sudo docker run -d --name $(NAME) -ti $(IMAGE):$(VERSION) /bin/bash

rm:
	-sudo docker stop $(NAME)
	sudo docker rm $(NAME)
	
