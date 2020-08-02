FROM golang:1.13

WORKDIR /go/src/github.com/fangbinwei/aliyun-oss-website-action
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

# run in /github/workspace
CMD ["aliyun-oss-website-action"]