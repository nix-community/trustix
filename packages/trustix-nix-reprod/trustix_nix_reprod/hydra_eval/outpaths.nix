{ systems }:
let
  path = <nixpkgs>;

  lib = import (path + "/lib");
  hydraJobs = import (path + "/pkgs/top-level/release.nix")
    # Compromise: accuracy vs. resources needed for evaluation.
    {
      supportedSystems = systems;

      nixpkgsArgs = {
        config = {
          allowBroken = false;
          allowUnfree = false;
          allowInsecurePredicate = x: true;
          checkMeta = false;

          handleEvalIssue = reason: errormsg:
            let
              fatalErrors = [
                "unknown-meta" "broken-outputs"
              ];
            in if builtins.elem reason fatalErrors
              then abort errormsg
              else true;

          inHydra = true;
        };
      };
    };
  recurseIntoAttrs = attrs: attrs // { recurseForDerivations = true; };

  # hydraJobs leaves recurseForDerivations as empty attrmaps;
  # that would break nix-env and we also need to recurse everywhere.
  tweak = lib.mapAttrs
    (name: val:
      if name == "recurseForDerivations" then true
      else if lib.isAttrs val && val.type or null != "derivation"
              then recurseIntoAttrs (tweak val)
      else val
    );

  # Some of these contain explicit references to platform(s) we want to avoid;
  # some even (transitively) depend on ~/.nixpkgs/config.nix (!)
  blacklist = [
    "tarball" "metrics" "manual"
    "darwin-tested" "unstable" "stdenvBootstrapTools"
    "moduleSystem" "lib-tests" # these just confuse the output
  ];

  system = "x86_64-linux";
in
  # { hello.${system} = (import path { inherit system; }).hello; }
  tweak (builtins.removeAttrs hydraJobs blacklist)
