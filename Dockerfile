FROM golang:alpine as build-env
RUN apk add gcc libc-dev
RUN apk add --update git
RUN apk add ca-certificates wget && update-ca-certificates
RUN git config --global url."https://Lawleeess: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@github.com/".insteadOf "https://github.com/"
ENV GO111MODULE=on
ADD . /go/src/github.com/Lawleeess/gympro-api
WORKDIR /go/src/github.com/Lawleeess/gympro-api
RUN go build -o gympro-api

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN apk --no-cache add tzdata
COPY --from=build-env /go/src/github.com/Lawleeess/gympro-api/gympro-api /go/gympro-api/gympro-api
COPY --from=build-env /go/src/github.com/Lawleeess/gympro-api/app.env /app.env

ENTRYPOINT ["/go/gympro/gympro-api"]