## Project structure

Trustix is structured as a monorepo consisting of many subpackages:

- [trustix](../../packages/trustix)

The main package with all log functionality.
This component is generic and doesn't know anything about any Nix or other package manager specifics.

- [trustix-doc](../../packages/trustix-doc)

The main documentation package that aggregates documentation from the various subpackages.

- [trustix-nix](../../packages/trustix-nix)

This is a supplemental daemon to the main Trustix daemon that layers some knowledge about Nix on top of the generic log functionality.
It contains a [post-build hook](https://www.tweag.io/blog/2019-11-21-untrusted-ci/) used to submit newly built packages to the logs, a binary cache HTTP interface and a development tool to submit already built closures.

- [trustix-nix-reprod](../../packages/trustix-nix-reprod)

This packages

- [trustix-proto](../../packages/trustix-proto)

Trustix-proto contains all shared protobuf definitions shared by various components, as well as generated Go libraries to interact with Trustix over it's RPC mechanism (gRPC).

- [trustix-python](../../packages/trustix-python)

Trustix-python contains generated code from trustix-proto for Python.
If you want to interact with Trustix over it's RPC interface from Python this is what you want to use.

- [pynix](./packages/pynix)

A number of generic small utility functions to work with Nix files.
At the time of writing this document it has an implementation of the Nix base32 encoding and a derivation file parser.

## Globally installed tooling

Trustix doesn't depend on much in the way of globally installed tools.

We do make two assumptions in regards to tooling managed outside of the repository though:

- [Nix](https://nixos.org)

If you've read this far you likely already know Nix and what it is, so we won't go into any detail about this.

- [direnv](https://direnv.net)

A shell extension to load directory local environments in a currently running shell and/or editor.
This will load a present `shell.nix`/`default.nix` when used with the direnv rule `use nix`, which is the mode of operation we are using direnv in.

## Getting started

All subpackages have their own shell environments which all needs to be explicitly whitelisted to be loaded.
For convenience we have a Makefile target in the root of the project called `direnv-allow`.
To whitelist all subpackages run:
``` sh
$ make direnv-allow
```

## Makefile structure

All components are using Makefile's as their development entry points for ease of use.

All standard Make targets are always implemented, even though they are no-ops in some cases where they don't make sense.
For example a build step doesn't make sense for most Python code.

These are all standard make targets you can expect to find for any given package:

- build

Builds the package.

- test

Runs the tests for a given package.

- lint

This target runs all configured linter steps.

- format

This target checks the formatting of a given package.

- develop

This target runs the package in development (watch) mode.

- doc

This target builds documentation.
This is mostly outputing markdown files in the relevant location for the `trustix-doc` package to compose.

## Running the whole setup

To run individual components change directory to the relevant package and run:
``` sh
$ make develop
```

This also works from the project root where it will start _all_ packages in watch mode.

## Quickly runing all tests

From the root directory run:
``` sh
$ make all
```

## Notes

Cryptographic keys for development is checked in to the repository for ease of use and a very quick getting started experience.
