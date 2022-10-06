# About

## [Trustix](https://github.com/nix-community/trustix) - build transparency reference implementation

[![asciicast](https://asciinema.org/a/384552.png)](https://asciinema.org/a/384552)

## Agree on inputs; agree on outputs

Trustix works by associating a hash of the inputs to a build with a hash of the resulting output. Reproducibility means that if multiple builders started with the same inputs, they should produce the same output. By comparing across builders, we can detect problems such as non-reproducibility and trojan viruses. By tracking the history of a build, we can ascertain security compromises, track build quality, and help build trust over time for a given system and provider.

## Binary planting protection

When we download a program, we typically trust that the binaries correspond to the program that we want. But how can we be sure that what we get is actually what we want, and not something malicious?

The trust we have in what we download is generally a proxy for the trust we have in the provider we download it from. This provider can put in place various measures to protect against attacks such as [man-in-the-middle attacks](https://en.wikipedia.org/w/index.php?title=Man-in-the-middle_attack), usually by using cryptographic hashes or digital signatures. These would protect against tampering in transit, but would not protect us if the binary was already compromised.

For example, the build infrastructure of the provider might be compromised, such that even if good source code goes in, bad binaries come out. This turns the provider into a single point of failure for trust.

Trustix makes it easy to distribute this responsibility among a number of non-authoritative peers. By comparing the build results between various providers, we can ascertain whether one or more might have been compromised.

## Tamper-evident history

Even if a provider's private key is compromised and an attacker gains unrestricted access, Trustix's immutable build log ensures that any modification to old entries can be detected and flagged. Any entries known to be from before the attack can still be trusted.

Since the build log provided by trustix is additive, we can ascertain when exactly a provider might have been compromised. Any binaries from before this point can still be trusted. Any modification of the log history is easily detected by peers, making history tampering impossible.

## Decide your level of trust

There is no universal notion of "verified"; you decide which providers to trust and how much. Your personal security needs may be low, in which case you can trust any binary cache on your list. In your high-security environment, however, you may require every provider to agree on the same output in order to trust any of them.


## Analyze package reproducibility

In the real world, builds are often non-reproducible without being maliciously modified. The build output might contain the time or date, or a reference to a local folder used during the build process. While most of these cases are benign, we can't really tell. The answer is to measure and track reproducibility. By comparing what different providers build on different machines and in different circumstances, we can track which packages are reproducible, and thus more trustworthy.
