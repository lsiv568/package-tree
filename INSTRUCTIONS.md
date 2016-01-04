# package-tree: Bambambam's Coding Challenge

At Bambambam, we have found little to no correlation between doing well in brain-teasers and "whiteboard coding" and being a world-class software engineer. On the other hand, we believe that having a candidate working on a scenario as close to the real-world as possible is the best way to see if somebody is a good fit for our team.

That's why we are sending you this coding challenge. Our goal is to have you write code and supporting artifacts that reflect the way you think and act about code in your professional life, not put you under the gun writing code with a whiteboard pen with somebody second guessing each step and a wall clock telling you there's only few minutes left.

## The problem we'd like you to solve

For this fictional problem, we ask you to write a package indexer.

*Packages* are executables or libraries that can be installed in a system, often via a package manager such as apt, RPM, or Homebrew. Many packages use libraries that are also made available as packages, so usually a package will require other packages being installed before you can install it in your system.

The system you are going to write will be responsible for collecting metadata on what packages are currently installed in a system, and accept or reject new installations or uninstallations depending on dependencies being present or not.

It will be a socket server, listening for TCP connections on port 8080. Many clientes can connect to this port at the same time, and once they do connect they may send a message following this structure:

```
<command>|<package>|<dependencies>\n
```

Where
* `<command>` is mandatory, and is either `INSTALL`, `UNINSTALL`, or `QUERY`
* `<package>` is mandatory, the name of the package referred to by the command, e.g. `mysql`, `openssl`, `pkg-config`, `postgresql`, etc.
* `<dependencies>` is optional, and if present it will be a comma-delimited list of packages that need to be present before `<package>` is installed. e.g. `cmake,sphinx-doc,xz`
* The message always ends with the character `\n`

Here are some sample messages:
```
INSTALL|cloog|gmp,isl,pkg-config\n
INSTALL|ceylon|\n
UNINSTALL|cloog|\n
QUERY|cloog|\n
```

Once the message is sent, the client will wait for the server to return a response code `1\n` or `0\n`.

The `<command>`s are as follows:
* `INSTALL` means that the given package should be marked as installed. The server must return `1\n` if the package was installed or `0\n` if the package *could ot be installed because of a missing dependency that needs to be installed first*.
* `UNINSTALL` means that the given package should be removed. The server must return `1\n` if the package was uninstalled or `0\n` if the package *could not be uninstalled because some other package depends on it*.
* `QUERY` means that the client wants to know if a given package is currently installed. The server should return `1\n` if the package is currently installed or `0\n` if it isn't.

### Technology choices and constraints
Although code at Bambambam is mostly written in Go and Ruby, you should feel free to write your solution in any language you prefer. Although we use and write libraries at Bambambam, in this specific exercise we would like to see as much of you own code as possible, so we ask you **not to use any library apart from your chosen runtime's standard library**. Testing code and build tools might use libraries, but not production code.

We would also ask you to write code that you woud consider production-ready, something you think other people would be happy supporting. Please don't forget to send us your automated tests and any other artifact needed to develop, build, or run your solution.

## The package we sent you

Together with this `INSTRUCTIONS.md` file, you should have recieved a tarball. In this tarball you will find:

* An executable file called `do-package-tree`, our test harness
* Another tarball, containing the Go source code for the executable mentioned above

### The test harness

This executable run an automated test suíte. We would like you to use this to verify your solution before sending it to us, and it can also be useful as functional tests during your development.

To run the test suíte, first make sure your server is up and listening on port `8080`. Then execute the following command:

```
$ ./package-tree-test
```

The tool will first tests for correctness, then try a robustness test. Both should pass before you submit your solution to the challenge, and once they both pass you will see a message like this:

```
================
All tests passed!
================
```

We have built several other features in the test suíte you might find helpful. To see them all, execute the following command:

```
$ ./package-tree-test -h
```

## Requirements

### Must Haves
These are the requirements your submission must fulfil to be considered correct.

* Send us all source code for test and production, and any artifacts (build scripts, etc.)
* Your code must pass the supplied test harness using different random seeds and concurrency factor up to 100
* Your code builds on the latest Ubuntu Docker image (if you need a runtime, e.g. JVM or Ruby, let us know which one)
* Your code is something you'd be confortable putting in production and having your team maintaining

### Should Have
These should be fulfilled, but if missing please write us a line on why.

* 

### Nice to Have
Stretch goals. If you fulfil these requirements you get bonus points, but they aren't required.

* Source control history (e.g. the `.git` directory or a link to Github)
* Design rationale
