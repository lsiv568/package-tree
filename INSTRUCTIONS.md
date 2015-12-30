# Package Tree, Bambambam's Coding Challenge

At Bambambam, we assess our engineering candidates on potential, skills, and culture fit.

It became somehow the norm in our industry to use brain teasers and whiteboards to assess a candidate's potential and skill. We don't believe in this format. Solving puzzles during interviews doesn't reflect the environment that an engineer would face once they've joined us, and we've found that the process penalises more introverted candidates.

Instead, we like people applying for engineering positions at Bambambam to demonstrate their abilities working in a very small but somewhat realistic proposed problem.

Our goal is to have you write code and supporting artifacts that reflect the way you think and act about code in your professional life.

## Problem description

For this fictional problem, we ask you to write a package indexer.

Packages are software artifacts that can be installed in a system, often via a package manager such as apt, rpm, or Homebrew. Packages offer some new feature to a system. These days so many packages have libraries in common, so a package will often require other packages to be installed before you can install it in your system.

The system you are going to write will be responsible for collecting metadata on what ackages are installed in a system and accept or reject new installations or uninstallations depending if dependencies are fulfilled.

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

Once the message was sent, the client will wait for the server to return a response code `1\n` or `0\n`.

The `<command>`s are as follows:
* `INSTALL` means that the given package should be marked as installed. The server must return `1\n` if the package was installed or `0\n` if the package *could ot be installed because of a missing dependency that needs to be installed first*.
* `UNINSTALL` means that the given package should be removed. The server must return `1\n` if the package was uninstalled or `0\n` if the package *could not be uninstalled because some other package depends on it*.
* `QUERY` means that the client wants to know if a given package is currently installed. The server should return `1\n` if the package is currently installed or `0\n` if it isn't.

## Technology choices and constraints
Although code at Bambambam is mostly written in Go and Ruby, you should feel free to write your solution in any language you prefer. Although we use and write libraries at Bambambam, in this specific exercise we would like to see as much of you own code as possible, so we ask you **not to use any library apart from your chosen runtime's standard library**. Testing code and build tools might use libraries, but not production code.

We would also ask you to write code that you woud consider production-ready, something you think other people would be happy supporting. Please don't forget to send us your automated tests and any other artifact needed to develop, build, or run your solution.

## The test harness

Together with this instructions file you should have received an executable file called `package-tree-test`. This is an automated test you must use to verify your program before sending it to us. First start your server and make sure it opens a server socket on port 8080, then run the test code:
```
$ ./package-tree-test
```

The tool will first tests for correctness, then try a robustness test. Both should pass before you submit your solution to the challenge, and once they both pass you will see a message like this:
```
================
All tests passed!
================
```
