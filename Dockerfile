FROM golang AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go .

RUN go build -o /hosts

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /hosts /hosts

COPY hosts.db hosts.db

EXPOSE 8080

ENTRYPOINT ["/hosts"]
