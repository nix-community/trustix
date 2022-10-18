# This file is used by Colmena to deploy trustix.dev and associated services

{
  meta = {
    nixpkgs = import ../pkgs.nix { };
  };

  defaults = { pkgs, ... }: {
    imports = [ ../nixos ];
    security.acme.acceptTerms = true;
    security.acme.defaults.email = "adisbladis@gmail.com";

    services.nginx = {
      recommendedGzipSettings = true;
      recommendedOptimisation = true;
      recommendedProxySettings = true;
      recommendedTlsSettings = true;
    };
  };

  "trustix.dev" = { name, nodes, ... }: {
    imports = [ ./trustix.dev ];
  };
}
