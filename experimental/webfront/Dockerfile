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
ADD . /go/src/github.com/rarenivar/project5799/webfront

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

