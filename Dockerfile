FROM scratch
MAINTAINER Peter Idah <Peter.Idah@gmail.com>
ADD . /
CMD ["/go-test-app"]
