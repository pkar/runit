# runit

Run applications, and stuff, watch for changes optionally and restart.
Main use case might be a lighter upstart/systemd/supervisord replacement
for use in Docker containers.

### Binaries

Can be found in [releases](https://github.com/pkar/runit/releases)

```bash
$ curl -o runit-v0.0.4.linux.tar.gz -L https://github.com/pkar/runit/releases/download/v0.0.4/runit-v0.0.4.linux.tar.gz
$ tar -xzvf runit-v0.0.4.linux.tar.gz
$ chmod +x runit && mv runit /usr/local/bin/
```

### Running

```bash
$ runit
  -alive
    	try to keep the command alive if it dies, you would use this for long running services like a server *optional
  -cmd string
    	command to run *required
  -loglevel int
    	logging level 1 is info (default 1)
  -wait
    	used with watch, this will wait for file changes and then run the cmd given *optional
  -watch string
    	path to directory or file to watch and restart cmd, the command will be run on startup unless wait is specified *optional

$ runit --cmd="echo blah" --watch=./
2015/08/04 21:46:01 running echo blah
blah
2015/08/04 21:46:01 captured child exited continue...
2015/08/04 21:46:05 event:  "foo": CREATE
2015/08/04 21:46:05 Detected new file foo
2015/08/04 21:46:05 restart event
2015/08/04 21:46:05 restarting
2015/08/04 21:46:05 killing subprocess
2015/08/04 21:46:05 running echo blah
blah
2015/08/04 21:46:05 captured child exited continue...
2015/08/04 21:46:05 event:  "foo": CHMOD
2015/08/04 21:46:05 event:  "foo": CHMOD

$ # long running processes with restart and watch
$ runit --restart --watch . --cmd="test/test.sh"
INFO 2015/02/03 20:54:59 runit.go:100: running test/test.sh
foo
foo
foo
foo
^C

$ # long running processes without watch
$ # process can be restarted by sending sighup to runit or
$ # or killing the subprocess cmd
$ # kill -SIGHUP $PID
$ runit --restart --cmd="test/test.sh"
2015/08/04 21:47:14 running test/test.sh
2015/08/04 21:47:14 running test/test.sh
foo
foo
foo
foo
foo
foo
foo
foo
foo
foo
^C2015/08/04 21:47:24 captured interrupt
```

### Development

1. Set gopath and go get github.com/pkar/runit
2. cd src/github.com/pkar/runit
3. ./make.sh <command>

### Tests

```bash
$ ./make.sh test
+ eval test
++ test
++ go test -cover .
ok  	github.com/pkar/runit	5.536s	coverage: 89.2% of statements
++ golint .
++ go tool vet --composites=false .
```
