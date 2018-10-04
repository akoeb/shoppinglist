# shopping list #

## description ##

This is a simple web application to maintain a shopping list. You can view and modify this list from your computer at home or from your cell phone in the store. Or let your partner add things to the list while you are running to get the stuff ;-)

## Backend ##

The backend is written in golang with the echo framework.

Compile the backend with:

* `go get` - to get all go dependencies
* `go build`- to build the application

## Frontend ##

Frontend is a vue.js application that uses the foundation css framework. Vue and foundation dependencies for the app itself are already in the /public folder, you only need to run `foundation build` to get the file public/css/app.css compiled out of the scss dir.

## docker image ##

The actual complete and runnable application is also packaged as docker image, run it with 

```bash
docker run -d --name CONTAINERNAME -v PATH_TO_SQLITE.db:/data/shoppinglist.db akoeb/shoppinglist
```

## development ##

If you want to develop on this application, you will need to have golang installed, with correct GOPATH, for the backend. To work on the frontend, you need to have foundation installed.

* Start watching for changes in your css files with `foundation watch`
* build the backend with `go build`
* run the backend with `./shoppinglist`
* open a browser, pointing to http://127.0.0.1:8080

## MAINTAINERS ##

* https://github.com/coatla
* https://github.com/akoeb

## License ##

GPLv3, see LICENSE file in this repo.