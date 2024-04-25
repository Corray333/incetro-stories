# Dockerfile.api
FROM golang:1.22.2-alpine3.18

WORKDIR /app

COPY . .
RUN apk add bash make musl-dev 
# RUN curl -fsSL \
#     https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
#     sh
# RUN go install github.com/swaggo/swag/cmd/swag@latest


# RUN go mod download

# RUN make goose-up