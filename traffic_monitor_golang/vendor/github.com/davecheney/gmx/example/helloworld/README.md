# helloworld example

Install the gmxc client

	go get github.com/davecheney/gmx/gmxc

Install this example 

	go get github.com/davecheney/gmx/example/helloworld

Run the example in the background

	$GOBIN/helloworld &

Query it via gmxc

	$GOBIN/gmxc -p $(pgrep helloworld) hello
	hello: world

