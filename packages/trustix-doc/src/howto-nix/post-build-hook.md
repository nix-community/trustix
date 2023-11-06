# Trustix - Usage via Nix

In the previous chapter we set up the main Trustix daemon. It's now time to actually start using it to publish build results.

## Requisites
- A NixOS installation using Flakes
- The basic setup from the previous chapter

## Setup
- `configuration.nix`:
``` nix
{{#include ../../../../examples/02_post_build/configuration.nix}}
```

## Effect
This sets up Nix with a post build hook that publishes any builds performed locally to your locally running log.
