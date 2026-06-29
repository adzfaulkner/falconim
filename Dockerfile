FROM golang:1.26.4-alpine as go

WORKDIR /go/src/app
COPY ./* ./

#RUN go get -d -v ./...
#RUN go install -v ./...

CMD ["app"]

FROM amaysim/serverless:4.14.4 as serverless

RUN mkdir /app

WORKDIR /app

COPY ./package.* ./

RUN npm install

