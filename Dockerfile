FROM golang:latest
# Copy the local package files to the container's workspace.
# RUN go get github.com/team142/gg-mapgen
ADD . /go/src/github.com/team142/gg-mapgen
# RUN cd ./src/github.com/team142/gg-mapgen
# RUN pwd

RUN go get ./src/github.com/team142/gg-mapgen
RUN go install ./src/github.com/team142/gg-mapgen
ENTRYPOINT /go/bin/gg-mapgen