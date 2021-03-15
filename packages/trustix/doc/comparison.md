---
title: "Trustix - A comparison of package security & Trustix use cases"
numbersections: true
author:
- Adam Höse
- Tweag I/O
email:
- adam.hose@tweag.io
lang: en-GB
classoption:
- twocolumn
header-includes:
- \usepackage{caption, graphicx, tikz, aeguill, pdflscape}
filters:
- pandoc-citeproc
bibliography:
- references.bib
---

# Introduction

## Comparison of package security models

### Classical signatures (NixOS/Debian/Arch Et al.)

The classical method for [code signing](https://en.wikipedia.org/wiki/Code_signing) is roughly:
1. Build a package, sometimes with declarative inputs but mostly not.
2. Create a hash over the binary contents produced, typically using sha256.
3. Create a signature over the package contents using a [public/private key pair](https://en.wikipedia.org/wiki/Public-key_cryptography).
4. Distribute package + signature.

This model has mostly served Linux distributions well over the years, though there has been documented cases where compromised build servers has been identified well after the fact.

Strengths:
- Very simple, easily understood
- Tiny overhead in packaging size
- Clients typically only need to know a single key

Weaknesses:
- Large degree of trust in central entities
- If relying on [PKI](https://en.wikipedia.org/wiki/Public_key_infrastructure), can have large complexity
- No guarantee source code maps to distributed binary
- A single key compromise takes down the entire model

### Binary transparency from Mozilla

This scheme piggybacks on top of [certificate transparency](https://tools.ietf.org/html/rfc6962) by publishing standard X509 certificates in a Certificate Transparency log.

Strengths:
- Very little custom code

Weaknesses:
- Can only verify that a binary was produced by Mozilla, not that it maps to any one version of the source code
- Relatively expensive lookups

### Debian build transparency

Very similar

### Nix using Trustix

The purely functional software deployment model is ideal for reproducability & independent verification of software builds.
Builds are described with an _exact_ set of specified inputs, making it much more likely that those inputs will map to specific outputs.

Compared to most other package managers Nix enforces hashing of every input to a build process and sandboxes the build.
Ff even a single hash at the root of the dependency graph changes it propagates to the rest of the graph.

However, until now there has been no way to efficiently verify that an input graph results in a set of deterministic outputs, other than trying to build according to the same build specification on your own machine to verify.
There can be impurities introduced by build-time timestamps, random paths created at build-time, CPU instruction set, timings and so forth.

By adopting a log-based approach we can completely replace the old simple signature based model and allow for

All of this also potentially applies to [GNU Guix](https://guix.gnu.org/), were the to adopt Trustix as the packaging model is the same as Nix.

Strengths:
- Can be used to independently verify results

Weaknesses:
- Much higher degree of complexity than previous solutions
- Has storage & network overhead costs
