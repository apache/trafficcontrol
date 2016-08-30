go build -ldflags "-X main.GitRevision=`git rev-parse HEAD` -X main.BuildTimestamp=`date +'%Y-%M-%dT%H:%M:%S'`"
