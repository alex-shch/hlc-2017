FROM golang

RUN apt-get update && apt-get install zip unzip

ENV PATH=${PATH}:/usr/local/go/bin GOROOT=/usr/local/go GOPATH=/go DOCKER_RUN=yes

RUN mkdir -p /go/src/github.com/alex-shch/hlc-2017

ADD . /go/src/github.com/alex-shch/hlc-2017

RUN go build github.com/alex-shch/hlc-2017 && go install github.com/alex-shch/hlc-2017

EXPOSE 80

CMD mkdir data && unzip /tmp/data/data.zip -d data/ >> /dev/null && /go/bin/hlc-2017

# docker build -t travels-go .
# docker run --rm -p 8080:80 -v "$PWD/testdata":/tmp/data travels-go
# docker tag travels-go stor.highloadcup.ru/travels/alive_yak
# docker push stor.highloadcup.ru/travels/alive_yak
