IMAGE := akoeb/shoppinglist
NAME := shoppinglist
VERSION := latest
HTTP_USER := my_list_user
HTTP_PW := supersecretpassword
#NETWORK := --network wordpress

.PHONY: all run bash rm

all: runlocal

compile:
	go build .

# migration uses https://github.com/golang-migrate/migrate
migratelocal:
	migrate -source file://db/ -database sqlite3://shoppinglist.db up

runlocal: compile
	foundation watch &
	./shoppinglist -db ./shoppinglist.db

build:
	docker build . -t $(IMAGE):$(VERSION)

run:
	docker run -d --name $(NAME) -v $(PWD)/shoppinglist.db:/data/shoppinglist.db -e HTTP_USER="$(HTTP_USER)" -e HTTP_PASSWORD="$(HTTP_PW)" -p 8080:8080 $(NETWORK) $(IMAGE):$(VERSION)

migrate:
	docker run --name $(NAME) -ti -v $(PWD)/shoppinglist.db:/data/shoppinglist.db $(IMAGE):$(VERSION) migrate -source file:///app/db/ -database sqlite3:///data/shoppinglist.db up

bash:
	docker run --name $(NAME) -ti $(IMAGE):$(VERSION) /bin/bash

rm:
	-docker stop $(NAME)
	docker rm $(NAME)
	
