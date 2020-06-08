FROM alpine:3.12.0

RUN apk add --no-cache libc6-compat

WORKDIR /app
COPY . .

WORKDIR doc
RUN chmod +x ./doc
ENTRYPOINT ["./doc"]

EXPOSE 8080