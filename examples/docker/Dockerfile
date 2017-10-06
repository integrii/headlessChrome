FROM yukinying/chrome-headless-browser
MAINTAINER integrii@gmail.com
USER root
RUN apt-get update
RUN apt-get -y install golang-1.8
RUN apt-get -y install git
RUN mkdir /go
RUN mkdir /app
ADD main.go /app/main.go
WORKDIR /app
RUN /usr/lib/go-1.8/bin/go get github.com/integrii/headlessChrome
RUN /usr/lib/go-1.8/bin/go build -o /app/headless
ENTRYPOINT ["/app/headless","--chromePath=/opt/google/chrome-unstable/chrome"]
