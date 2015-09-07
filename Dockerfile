FROM golang:onbuild

EXPOSE 8050

CMD ["/go/bin/app"]
