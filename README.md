# An Insecure Go Playground for the Adventurers

![Go Playground](https://assets.xiam.io/projects/go-playground.png)

The official [Go Playground][2] runs on a [sandbox][5], which is an intelligent
way of executing code from untrusted sources since it minimizes security risks.

However, if you ever want to showcase features that require network access, a
real filesystem, CGO, or anything that is not supported or actively restricted
by the sandbox you're out of luck.

This **unrestricted and insecure playground** offers you the sense of adventure
and danger you're looking for by removing all restrictions and security
features while keeping API compatibility with the official Go sandbox.

**Do not use this playground** unless you're sure you know what you're doing,
taking enough security measures on a different layer, and you know that no
matter what you do, people will try to abuse the system and probably succeed.

If you made it here, it means your sense of adventure is ticking; warnings
aside, feel free to use this playground in your next workshop to demonstrate
your Go projects to others.

## Running with docker

There are two different components to run this playground: one is the web
interface, and the other is the executor. The web interface is simple
application for the web that connects to an executor and displays a result. The
executor is the component that compiles, executes, and stores the output of a
Go program.

Executing untrusted code is risky enough; the recommended way of running the
playground is by isolating it in a [Docker][4] container and taking enough
measures to control and restrict it to suit your needs.

To run the web interface using Docker, use this command:

```sh
docker run \
    --rm \
    --name go-playground-webapp \
    -p 127.0.0.1:3000:3000 \
    xiam/go-playground:latest \
        go-playground-webapp
```

The example above will connect by default to the official sandboxed executor.
You'll need to run your own executor to run unrestricted code snippets.

To run the unresticted executor using Docker, use this command:

```sh
docker run \
    --rm \
    --name go-playground-executor \
    -p 0.0.0.0:3003:3003 \
    xiam/go-playground:latest \
        go-playground-executor
```

To make the web interface use the custom executor, you'll have to connect both
containers to the same [docker network][6], and specify the address of the
executor using the `-c` flag.

```sh
docker run \
    --rm \
    --network go-playground \
    --name go-playground-webapp \
    -p 127.0.0.1:3000:3000 \
    xiam/go-playground:latest \
        go-playground-webapp \
            -c https://go-playground-executor:3003
```

# License

This project is derived from the official [Go Playground][3] and is licensed
under the BSD license. See the [LICENSE](LICENSE) file for details.

[1]: https://www.golang.org/
[2]: https://play.golang.org/
[3]: https://github.com/golang/playground
[4]: https://www.docker.com/
[5]: https://en.wikipedia.org/wiki/Sandbox_(computer_security)
[6]: https://docs.docker.com/engine/network/
