#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at
#
#        http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS,
#    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#    See the License for the specific language governing permissions and
#    limitations under the License.
#

# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
# Note: need to move to ~/go?
ADD . /go/src/github.com/rarenivar/project5799/auth

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/jmoiron/sqlx
RUN go get github.com/lib/pq
WORKDIR ./src/github.com/rarenivar/project5799/auth
RUN export GOBIN=$GOPATH/bin && go install auth.go

# Document that the service listens on port 9000.
EXPOSE 8080

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/auth auth.config

