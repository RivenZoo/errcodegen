# Version: 1.0
FROM golang:1.12 as builder

RUN mkdir -p /data/app
ADD ./ /data/app
RUN cd /data/app && make clean build

FROM golang:1.12

RUN mkdir -p /data/app

WORKDIR /data/app

COPY --from=builder /data/app/build /usr/local/bin

