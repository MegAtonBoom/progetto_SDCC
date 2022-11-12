FROM golang:1.17-alpine
ARG service
ENV ser=${service}
ADD . /mutual_exclusion
WORKDIR /mutual_exclusion
COPY . .
RUN go mod download
RUN go build
ENTRYPOINT ./main ${ser}
