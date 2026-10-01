FROM golang:1.26

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# any version other than Dev also makes clover serve its embedded assets rather than reading ./static off disk
ARG VERSION=Dev
ARG DATE=unknown

COPY . .
RUN go build -v -ldflags "-X main.version=${VERSION} -X main.date=${DATE}" -o /usr/local/bin/app github.com/nyaruka/clover/v26/cmd/clover

EXPOSE 8060
CMD ["app"]
