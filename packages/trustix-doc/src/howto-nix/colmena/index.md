# Trustix - Builder setup with Colmena (flake based)

Up until now we have talked about components in isolation, let's take a stab at a real deployment using [Colmena](https://github.com/zhaofengli/colmena).
Once we're done we will have:
- A Trustix instance that others can subscribe to
- A local post build hook that publishes builds

Most of the contents in this article will be applicable to other deployment systems too with only minimal changes, the biggest difference should be how they deal with key material (secrets).

## Requisites / assumptions

- You need a recent Nix with Flakes enabled
- We start off in an empty git repository
- Our domain name is `demo.trustix.dev`

## Create keys

All Trustix build logs are first and foremost identified by their key pair, which will be the first thing we have to generate.

Let's start by generating a key pair for our log:
```
$ mkdir secrets
$ nix run github:nix-community/trustix#trustix -- generate-key --privkey secrets/log-priv --pubkey secrets/log-pub
```

Additionally logs are identified not just by their key, but how that key is used.
If a key is used for multiple protocols (not just Nix) those logs will have a different ID.
This ID is what _subscribers_ use to indicate what they want to subscribe to.

To find out the log ID for the key pair you just generated:
`$ nix run github:nix-community/trustix#trustix -- print-log-id --protocol nix --pubkey $(cat secrets/log-pub)`

## Create a deployment

In your `flake.nix` put:
``` nix
{{#include ../colmena/flake.nix}}
```

Enter the development shell and deploy:
```
$ nix develop  # Pulls in Colmena via Flakes devShells
$ colmena apply  # Deploy
```

## Spread your log
In the next chapter we will go over how to use this log from clients.
The most important thing right now is to make a note of your _public_ key and your _domain name_.
