# Builder
FROM golang:1.19.4-alpine3.17 as builder

RUN apk update && apk upgrade && \
	apk --update add git make bash build-base

WORKDIR /app

COPY . .

RUN go build -trimpath -o bcrypt_cost_factor_benchmark main.go
RUN ls

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
	apk --update --no-cache add tzdata

COPY --from=builder /app/bcrypt_cost_factor_benchmark /

ENV PORT=3000

EXPOSE 3000

CMD ["/bcrypt_cost_factor_benchmark"]