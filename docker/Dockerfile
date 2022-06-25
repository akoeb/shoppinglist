# STEP 1 build executable binary
FROM golang:1.16-alpine as gobuilder

# Install SSL ca certificates
RUN apk update && apk add ca-certificates build-base

# Create appuser
RUN adduser -D -g '' appuser
COPY backend $GOPATH/src/github.com/akoeb/shoppinglist/
WORKDIR $GOPATH/src/github.com/akoeb/shoppinglist/

# get dependancies for go
RUN go get -d -v 

# build the binary
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -tags netgo -ldflags='-w -extldflags "-static"' -o /go/bin/shoppinglist


# STEP 2: frontend builder
FROM node:15-alpine as nodebuilder

COPY frontend /home/node

RUN chown -R node:node /home/node

USER node
WORKDIR /home/node

# get dependencies and install frontend
RUN  npm install && npm run build-prod


# STEP 3 build a small image
# start from alpine to have a shell in the image, for the entry point
FROM alpine
COPY --from=gobuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=gobuilder /etc/passwd /etc/passwd

# Copy our static executable
COPY --from=gobuilder /go/bin/shoppinglist /shoppinglist
COPY --from=nodebuilder /home/node/dist /public
ENV DATABASE_FILE="/data/shoppinglist.db"
ENV DOMAIN="localhost"
ENV HTTP_USER=""
ENV HTTP_PASSWORD=""
USER appuser
VOLUME /data
CMD /shoppinglist -db /data/shoppinglist.db -domain $DOMAIN -user $HTTP_USER -password $HTTP_PASSWORD