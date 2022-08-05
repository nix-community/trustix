[![CI](https://github.com/tweag/trustix/actions/workflows/ci.yml/badge.svg)](https://github.com/tweag/trustix/actions/workflows/ci.yml)

# Trustix - Distribute trust
Trustix aims to pave the way for an internet where small providers can deliver safe
code, ultimately increasing security and number of trusted packages for all users.

It does it by comparing results of software builds across independent, decentralized
providers in order to:

- _Assure trust_ in downloaded software binaries by automatically identifying
  and exclude misbehaving builders and their generated artifacts
- Track _build reproducibility_, making it easier to identify packages that
  do not reliably produce the same output when build

Trustix is a novel trust assurance model for the
[Nix ecosystem](https://nixos.org/) developed via an
[NGI0 PET grant](https://nlnet.nl/project/Trustix/).

## Overview
Pre-built software binaries downloaded from centralized providers are implicitely
trusted to correspond to the desired program. However, nothing guarantees that
these binaries were really built from the program's actual source code using
reasonable build instructions.

Costly supply chain attacks usually exploit this situation to distribute malicious
software, which is one reason why most software is delivered through centralized,
highly secured providers in the first place.

Trustix aims to solve this problem via distributed trust.

### Distributed trust
The essential idea behind Trustix is to compare build outputs across a group of
independent builders that log and exchange hashes of build input/output pairs.
This is achieved through the following methodology:

- Each builder is associated with a public-private key pair
- In a post-build hook the output hash (NAR hash) of the build is uploaded to a
  ledger (a signed append-only log of build results).
- The results from different builders are aggregated and compared by Trustix.

This allows a user to trust binary substitutions based on an M-of-N vote among the
participating builders.

### Example
Let's say there are 4 builders available: `mercury`, `venus`, `saturn` and `mars`.

All builders communicate precisely what they have built with a hash that describes
the build inputs of a given package, and another that describes the output obtained.
They all claim to have built the `hello` derivation.

A Trustix instance is configured to require a 3/4 majority for a build to be considered
safe and is subscribed to the logs of all 4 builders above.

Upon inspection of the results of the 4 builders, Trustix observes that for the same
input, 3 builders have arrived at the same output hash but `saturn` has obtained
something different:

- mercury: `fe726118f9f6ecd9739554ac16f32b499ad7a981`
- venus:   `fe726118f9f6ecd9739554ac16f32b499ad7a981`
- saturn:  `a4f208705a0995a1dbe287916d44b38ad85e770b` <- _This hash is different from the others_
- mars:    `fe726118f9f6ecd9739554ac16f32b499ad7a981`

With this information Trustix could:

- Automatically identify and exclude misbehaving builders such as `Chuck` in
  the example above, preventing the installation of unsafe software.
- Trust only builds that have been confirmed by a majority of selected builders.
- Track build reproducibility across a large number of builds. Perhaps `saturn`
  is not a bad actor and the build process for a package is simply not reproducible.

## How does it relate to Nix?
In the Nix ecosystem, pre-built binaries are distributed through so-called
_binary substituters_. Similar to other centralized caching systems, they are a single
point of failure in the chain of trust when delivering a package to a user. This is
problematic for the following reasons:

- If anyone manages to compromise the [NixOS Hydra](https://hydra.nixos.org/) build
  machines and its keys, they could upload backdoored builds to users. In the Nix
  ecosystem, a compromised key is even more dangerous because https://cache.nixos.org
  can't use a rolling key because of the way it is set up. This means that a
  compromised key would realistically mean that _all_ packages in the cache are
  compromised. They would have to be rebuilt or garbage collected which is very costly.

- The NixOS Hydra _hardware_, on which the binaries are built, may also be compromised
  and not considered trustworthy by more security conscious users.

For more background information, see the original
[project announcement](https://www.tweag.io/blog/2020-12-16-trustix-announcement/).

## Documentation

Documentation is built as a part of CI and published on
[Github Pages](https://tweag.github.io/trustix/).

## Developing

For notes on development see [HACKING.md](./packages/trustix-doc/src/hacking.md)
