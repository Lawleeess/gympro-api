FROM golang:alpine as build-env
RUN apk add gcc libc-dev
RUN apk add --update git
RUN apk add ca-certificates wget && update-ca-certificates
RUN git config --global url."https://Lawless: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@github.com/".insteadOf "https://github.com/"
ENV GO111MODULE=on
ENV GOPRIVATE="github.com/Lawleeess/*"
ADD . /go/src/github.com/CValier/gympro-api
WORKDIR /go/src/github.com/CValier/gympro-api
RUN go build -o gympro-api

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN apk --no-cache add tzdata
COPY --from=build-env /go/src/github.com/CValier/gympro-api/gympro-api /go/gympro-api/gympro-api
COPY --from=build-env /go/src/github.com/CValier/gympro-api/app.env /app.env

ENTRYPOINT ["/go/gympro-api/gympro-api"]