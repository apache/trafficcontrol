# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
# Note: need to move to ~/go?
ADD . /go/src/github.com/rarenivar/project5799/webfront

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/dgrijalva/jwt-go
#RUN ls -lR
WORKDIR ./src/github.com/rarenivar/project5799/webfront
#RUN cd ./src/github.com/rarenivar/project5799/webfront && export GOBIN=$GOPATH/bin && go install webfront.go
RUN export GOBIN=$GOPATH/bin && go install webfront.go
#RUN ls -l
#RUN go install webfront.go

# Document that the service listens on port 9000.
EXPOSE 9000

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/webfront webfront.config

