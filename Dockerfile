FROM docker.io/library/golang:1.21

WORKDIR /build
COPY . .
RUN make

FROM gcr.io/distroless/static-debian12:latest

COPY --from=builder /build/build/ms1-server /

EXPOSE 50051

ENTRYPOINT ["/ms1-server"]
CMD ["run"]
