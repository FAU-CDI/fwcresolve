FROM docker.io/library/alpine:3.14 as os

# install ca-certificates
RUN apk add --update --no-cache ca-certificates

# create www-data
RUN set -x ; \
  addgroup -g 82 -S www-data ; \
  adduser -u 82 -D -S -G www-data www-data && exit 0 ; exit 1

FROM docker.io/library/golang:1.18 as builder
ADD . /app/
WORKDIR /app/
RUN go get ./cmd/wdresolve
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wdresolve ./cmd/wdresolve

# add it into a scratch image
FROM docker.io/library/alpine:3.14

# add the user
COPY --from=os /etc/passwd /etc/passwd
COPY --from=os /etc/group /etc/group

# grab ssl certs
COPY --from=os /etc/ssl/certs /etc/ssl/certs

# add the app
COPY --from=builder /app/wdresolve /wdresolve

# and set the entry command
EXPOSE 8080
USER www-data:www-data
CMD ["/wdresolve", "-listen", "0.0.0.0:8080"]