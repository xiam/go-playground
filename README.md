# Your Own Go Playground

The official [Go Playground][1] runs on a [sandbox][5], which is the
recommended way of doing things like executing code from untrusted sources,
since it minimizes security risks.

However, if you ever want to showcase features that require network access, a
real filesystem, CGO or any kind of thing that is not supported or actively
restricted by [NaCL](https://developer.chrome.com/native-client) you're out of
luck.

This *unrestricted and insecure* playground offers you more flexibility on that
front, you can choose when to use the sandboxed environment and when not to.

## Quick start with docker

You can run playground and connect it to a local sandbox like this:

```
docker run \
  -d \
  --name go-playground-sandbox \
  -p 127.0.0.1:8080:8080 \
  xiam/go-playground-sandbox

# Running unsafebox
# docker run \
#  -d \
#  --name go-playground-unsafebox \
#  -p 127.0.0.1:8080:8080 \
#  xiam/go-playground-unsafebox

# Running web editor
docker run \
  -d \
  --link go-playground-sandbox:compiler \
  --name go-playground \
  -p 0.0.0.0:3000:3000 \
  xiam/go-playground \
    bash -c \
      'webapp -c http://compiler:8080/compile?output=json'
```

## Front-end

This is similar to the
[play.golang.org](https://github.com/golang/playground/tree/master/app) web
app, except that it:

* Does not depend on [appengine](https://cloud.google.com/appengine/docs/go/reference).
* Uses [boltdb](https://github.com/boltdb/bolt) to save data.
* Can be configured to communicate with any other Go Playground service (the
  one that compiles and runs Go code), including the official one.

You can build and run a Go Playground like this:

```
cd webapp
go build
./webapp -allow-share
# 2019/01/26 20:35:47 Serving Go playground at :3000 (with compiler https://play.golang.org/compile?output=json)
```

This will create a local server that uses the official Go Playground service to
build and run Go code. See the main page at http://127.0.0.1:3000.

![screen shot 2016-01-03 at 8 12 25 am](https://cloud.githubusercontent.com/assets/385670/12079146/1de8c24a-b1f4-11e5-87b9-10f0a22054e5.png)

Love those nice [live
examples](https://golang.org/pkg/strings/#example_Contains) on the golang.org
site? You can also embed them in your website using a few lines of code:

![screen shot 2016-01-03 at 8 12 50 am](https://cloud.githubusercontent.com/assets/385670/12079219/9fd19f14-b1f6-11e5-949e-f36561a7f0ff.png)

See a local example at http://127.0.0.1:3000/example.

### Unsafebox

I basically took out all the security measures of the original sandbox and
generated a dumbed down version of it which actually does not sandbox anything
and will put your life at risk.

This is a unrestricted linux/amd64 installation, you should not really use this
box unless you're absolutely sure you know what you're doing and you're aware
that no matter what you do, people will try to abuse the system and they will
probably succeed.

Ok, now that you're warned, you can (but shouldn't) build and run this
*dangerous box* like this:

```
cd unsafebox
make docker-run
```

You can point the Go Playground web app to this service using the `-c`
parameter:

```
cd webapp
go build
./webapp -allow-share -c "http://localhost:8080/compile?output=json"
```

Remember that this machine is completely open to the world, if you plan to
upload it to a public place you should take some other containment measures,
such as not using root to run the sandbox, using chroot jails and iptables
rules, etc. this really depends on your specific needs.

### Importing custom packages

Users of your playground won't be able to install or use packages that are not
part of the Go standard library, in case you want to showcase a special package
you'll have to create a slightly different docker image on top of the sandbox
or the unsafebox, see this `Dockerfile`:

```
FROM xiam/go-playground-unsafebox

RUN go get github.com/myuser/mypackage
RUN go get github.com/otheruser/otherpackage
```

You can build that docker image and then start `webapp` using the `-c`
parameter pointing to the docker image and you'll be able to import custom
packages from your playground:

```
./webapp -c "http://custom.box/compile?output=json"
```

![screen shot 2016-01-03 at 2 32 00 pm](https://cloud.githubusercontent.com/assets/385670/12080650/d6037186-b226-11e5-8bd1-3b98627a1e03.png)

Feel free to use this playground on your next workshop to demonstrate your Go
projects to others.

[1]: https://www.golang.org/
[2]: https://play.golang.org/
[3]: https://github.com/golang/playground
[4]: https://www.docker.com/
[5]: https://en.wikipedia.org/wiki/Sandbox_(computer_security)
