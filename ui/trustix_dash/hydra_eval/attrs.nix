let
  pkgs = import ./. { };
  inherit (pkgs) lib;

  supportedSystems = [ "x86_64-linux" ];

  # Strip most of attributes when evaluating to spare memory usage
  scrubJobs = true;

  nixpkgsArgs = { config = { allowUnfree = false; inHydra = true; }; };

  rLib = import (pkgs.path + "/pkgs/top-level/release-lib.nix") { inherit supportedSystems scrubJobs nixpkgsArgs; };

in rLib.mapTestOn (rLib.packagePlatforms pkgs)
