# package-tree: Bambambam's Coding Challenge

At Bambambam, we have found little to no correlation between doing well in brain-teasers and "whiteboard coding" and being a world-class software engineer. On the other hand, we believe that having a candidate working on a scenario as close to the real-world as possible is the best way to see if somebody is a good fit for our team.

That's why we are sending you this coding challenge. Our goal is to have you write code and supporting artifacts that reflect the way you think and act about code in your professional life, not put you under the gun writing code with a whiteboard pen with somebody second-guessing each step and a wall clock telling you there are only a few minutes left.

## The problem we'd like you to solve

For this fictional problem, we ask you to write a package indexer.

*Packages* are executables or libraries that can be installed in a system, often via a package manager such as apt, RPM, or Homebrew. Many packages use libraries that are also made available as packages themselves, so usually a package will require you to install its dependencies you can install it on your system.

The system you are going to write keeps track of package dependencies. Clients will connect to your server and inform which packages should be indexed, and which dependencies they might have on other packages. We want to keep our index consistent, so your server must not index any package until all of its dependencies have been indexed first. The server should also not remove a package if any other packages depend on it.

The server will open a TCP socket on port 8080. It must accept connections from multiple clients at the same time, all trying to add and remove items to the index concurrently. Clients are independent of each other, and it is expected that they will send repeated or contracting messages. New clients can connect and disconnect at any moment, and sometimes clients can behave badly and try to send broken messages.

Messages from clients follow this pattern:

```
<command>|<package>|<dependencies>\n
```

Where:
* `<command>` is mandatory, and is either `INDEX`, `REMOVE`, or `QUERY`
* `<package>` is mandatory, the name of the package referred to by the command, e.g. `mysql`, `openssl`, `pkg-config`, `postgresql`, etc.
* `<dependencies>` is optional, and if present it will be a comma-delimited list of packages that need to be present before `<package>` is installed. e.g. `cmake,sphinx-doc,xz`
* The message always ends with the character `\n`

Here are some sample messages:
```
INDEX|cloog|gmp,isl,pkg-config\n
INDEX|ceylon|\n
REMOVE|cloog|\n
QUERY|cloog|\n
```

For each message sent, the client will wait for a response code from the server. Possible response codes are `OK\n`, `FAIL\n`, or `ERROR\n`. After receiving the response code, the client can send more messages.

The response code returned should be as follows:
* For `INDEX` commands, the server returns `OK\n` if the package could be indexed or if it was already present. It returns `FAIL\n` if the package cannot be indexed because some of its dependencies aren't indexed yet and need to be installed first.
* For `REMOVE` commands, the server returns `OK\n` if the package could be removed from the index. It returns `FAIL\n` if the package could not be removed from the index because some other indexed package depends on it.
* For `QUERY` commands, the server returns `OK\n` if the package is indexed. It returns `FAIL\n` if the package isn't indexed.
* If the server doesn't recognise the command or if there's any problem with the message sent by the client it should return `ERROR\n`.

### Technology choices and constraints
Code at Bambambam is mostly written in Go and Ruby, you should feel free to write your solution in any language you prefer. Although we use and write libraries at Bambambam, for this exercise we would like to see as much of you own code as possible, so we ask you **not to use any library apart from your chosen runtime's standard library**. Testing code and build tools might use libraries, but not production code.

We would also ask you to write code that you woud consider production-ready, something you think other people would be happy supporting. Please don't forget to send us your automated tests and any other artifact needed to develop, build, or run your solution.

## The package we sent you

Together with this `INSTRUCTIONS.md` file, you should have received a tarball. In this tarball you will find:

* A Linux executable file called `do-package-tree`, our test harness
* Another tarball, containing the Go source code for the executable mentioned above

### The test harness

This executable runs an automated test suite. We would like you to use this to verify your solution before sending it to us, and it can also be useful as functional tests during your development.

To run the test suite, first make sure your server is up and listening on port `8080`. Then execute the following command:

```
$ ./package-tree-test
```

The tool will first test for correctness, then try a robustness test. Both should pass before you submit your solution to the challenge, and once they both pass you will see a message like this:

```
================
All tests passed!
================
```

We have built several other features in the test suite you might find helpful. To see them all, execute the following command:

```
$ ./package-tree-test --help
```

## Requirements

### Must Haves
These are the requirements your submission must fulfil to be considered correct.

* Send us all source code for test and production, and any artifacts (build scripts, etc.)
* Your code is something you'd be confortable putting in production and having your team maintaining
* Your code must pass the supplied test harness using different random seeds and concurrency factor up to 100
* Your code builds on the latest Ubuntu Docker image (if you need a runtime, e.g. JVM or Ruby, let us know which one)

### Should Have
These should be fulfilled, but if they aren't please let us know why.

* Your code is tested in some automated fashion at unit and integration levels
* Source control history (e.g. the `.git` directory or a link to Github)

### Nice to Have
Stretch goals. If you fulfil these requirements, you get bonus points, but they aren't required.

* Design rationale
* A Dockerfile
