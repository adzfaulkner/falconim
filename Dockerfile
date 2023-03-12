FROM golang:1.20-alpine as go

WORKDIR /go/src/app
COPY ./* ./

#RUN go get -d -v ./...
#RUN go install -v ./...

CMD ["app"]

FROM amaysim/serverless:3.23.0 as serverless

RUN mkdir /app

WORKDIR /app

COPY ./package.* ./

RUN npm install

