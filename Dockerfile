FROM alpine
RUN apk update && apk --no-cache add tzdata
ENTRYPOINT ["/traindepart-mqtt"]
COPY traindepart-mqtt /