# STEP 1 build executable binary
FROM golang:alpine as builder

# Install SSL ca certificates
RUN apk update && apk add git && apk add ca-certificates && apk add build-base && apk add --update nodejs nodejs-npm && apk add python2

# Create appuser
RUN adduser -D -g '' appuser
COPY . $GOPATH/src/github.com/akoeb/shoppinglist/
WORKDIR $GOPATH/src/github.com/akoeb/shoppinglist/

# get dependancies for go
RUN go get -d -v 

# build the binary
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -tags netgo -ldflags='-w -extldflags "-static"' -o /go/bin/shoppinglist

# get dependencies and install frontend
RUN  npm install && node_modules/.bin/foundation build


# STEP 2 build a small image
# start from alpine to have a shell in the image, for the entry point
FROM alpine
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Copy our static executable
COPY --from=builder /go/bin/shoppinglist /shoppinglist
COPY --from=builder /go/src/github.com/akoeb/shoppinglist/public /public
ENV DATABASE_FILE="/data/shoppinglist.db"
ENV DOMAIN="localhost"
ENV HTTP_USER=""
ENV HTTP_PASSWORD=""
USER appuser
VOLUME /data
CMD /shoppinglist -db /data/shoppinglist.db -domain $DOMAIN -user $HTTP_USER -password $HTTP_PASSWORD