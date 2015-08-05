# runit

Run applications, and stuff, watch for changes optionally and restart.
Main use case might be a lighter upstart/systemd/supervisord replacement
for use in Docker containers.

### Binaries

Can be found in [releases](https://github.com/pkar/runit/releases)

```bash
$ curl -o runit-v0.0.2.linux.tar.gz -L https://github.com/pkar/runit/releases/download/v0.0.2/runit-v0.0.2.linux.tar.gz
$ tar -xzvf runit-v0.0.2.linux.tar.gz
$ chmod +x runit && mv runit /usr/local/bin/
```

### Running

```bash
$ runit -cmd="echo blah" -watch=./
INFO 2015/02/03 20:54:23 runit.go:100: running echo blah
blah
^C

$ # long running processes with keep alive and watch
$ runit -alive -cmd="test/test.sh" -watch=./
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
$ runit -alive -cmd="test/test.sh"
INFO 2015/02/03 20:54:59 runit.go:100: running test/test.sh
foo
foo
foo
foo
^C
```

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
