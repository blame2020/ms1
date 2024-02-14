FROM docker.io/library/golang:1.21 as builder

RUN pwd
COPY . /build
WORKDIR /build
RUN ls -al
RUN make

FROM gcr.io/distroless/static-debian12:latest

COPY --from=0 /build/build/ms1-server /

EXPOSE 50051

ENTRYPOINT ["/ms1-server"]
CMD ["run"]
