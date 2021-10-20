# Trustix - A new model for Nix binary substitutions

Trustix is a tool that compares build outputs for a given build input across a group of independent providers to establish trust in software binaries.

## Overview

We often use pre-built software binaries and trust that they correspond to the program we want.
But nothing assures that these binaries were really built from the program's sources and reasonable built instructions.
Common, costly supply chain attacks exploit this to distribute malicious software, which is one reason why most software is delivered through centralized, highly secured providers.
Trustix, a tool developed via an [NGI0 PET grant](https://nlnet.nl/project/Trustix/), establishes trust in binaries in a different, decentralized manner.
This increases security, and paves the way for an internet where small providers can deliver safe code, ultimately with a safer and larger offer for the user.

Trustix is developed for the Nix ecosystem.

## How does this translate to Nix?

In the Nix ecosystem, pre-built binaries are distributed through so-called _binary substituters_.
Similar to other centralized caching systems, they are a single point of failure in the chain of trust when delivering a package to a user.
This is problematic for several reasons:

First, if anyone manages to compromise the [NixOS Hydra](https://hydra.nixos.org/) build machines and its keys, they could upload backdoored builds to users.
In the Nix ecosystem, a compromised key is even more dangerous because https://cache.nixos.org can't use a rolling key because of the way it is set up.
This means that a compromised key would realistically mean that _all_ packages in the cache are compromised.
  They would have to be rebuilt or garbage collected which is very costly.

Second, the NixOS Hydra _hardware_, on which the binaries are built, may also be compromised and not considered trustworthy by more security conscious users.

For some more background see the original [project announcement](https://www.tweag.io/blog/2020-12-16-trustix-announcement/).

## Trustix design

`Trustix` aims to solve this problem via distributed trust & trust agility.
Essentially it compares build outputs across a group of independent builders
that log and exchange hashes of build input/output pairs.
This is achieved through the following methodology:

-   Each builder is associated with a public-private key pair
-   In a post-build hook the output hash (NAR hash) of the build is uploaded to a ledger (a signed append-only log of build results).

This allows a user to trust binary substitutions based on an M-of-N vote among the participating builders.

Here is an example:
Let's say we have 4 builders configured: `Alice`, `Bob`, `Chuck` & `Dan`.
We have configured `Trustix` to require a 3/4 majority for a build to be trusted.
`Alice`, `Bob`, `Dan` and `Chuck` all claim to have built the `hello` derivation.
All builders participate in the Trustix network and communicate precisely
what they have built with a hash that describes the build inputs of `hello`, and
what have obtained as output with another hash.
For the same input, the first 3 builders have arrived at the same output hash but `Chuck` has
obtained something different.

This information can now be used by a Trustix user to:

-   track build reproducability across a large number of builders.
-   trust only builds that have been confirmed by a majority of selected
    builders.
-   automatically identify and exclude misbehaving builders such as `Chuck` in
    above's example.

## Documentation

Documentation is built as a part of CI and published on [Github Pages](https://tweag.github.io/trustix/).

## Developing

For notes on development see [HACKING.md](./packages/trustix-doc/src/hacking.md)
