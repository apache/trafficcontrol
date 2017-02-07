# gmx client

`gmxc` is a simple command line client for interacting with gmx enabled processes.

## Usage

### Listing processes

The default invocation of `gmxc` will list the accessible gmx enabled processes currently running

	% ./gmxc 
	.gmx.16207.0    ["./godoc" "-v" "-http=:8080"]

### Retrieving gmx values

	./gmxc -p 16207 runtime.numcpu
	runtime.numcpu: 4

	./gmxc -p 14968 'runtime.(numcpu|version)'
	runtime.version: weekly.2012-01-27 11662
	runtime.numcpu: 2

### Listing all gmx values

	% ./gmxc -p 16207 keys
	keys: [keys runtime.memstats runtime.gomaxprocs runtime.version os.args runtime.numcpu runtime.cgocalls]

