FROM golang:1.26.4-alpine as go

WORKDIR /go/src/app
COPY ./* ./

#RUN go get -d -v ./...
#RUN go install -v ./...

CMD ["app"]

FROM node:26-alpine3.23 as serverless

RUN npm i serverless -g \
    && mkdir /app

WORKDIR /app

COPY ./package.* ./

RUN npm install

