{ config, pkgs, ... }:

{
  # Static website for trustix.dev
  services.nginx.virtualHosts."trustix.dev" = {
    root = pkgs.trustix-doc;
  };

  # Redirect old build-transparency.org domain to trustix.dev
  services.nginx.virtualHosts."build-transparency.org" = {
    extraConfig = ''
      return 302 https://trustix.dev$request_uri;
    '';
  };

  # Reproducibility dashboard
  services.trustix-nix-r13y = {
    enable = true;
    hostName = "r13y.trustix.dev";

    lognames = {
      e0f263745e4e3ab07ab5275b00b44f594e0b6d2bd35892a8ebd10a7f86322eb7 = "trustix-demo";
    };

    attrs = {
      nixos-unstable = [ "nixos.iso_minimal.x86_64-linux" ];
    };

    channels.hydra.nixos-unstable = {
      base_url = "https://hydra.nixos.org";
      jobset = "trunk-combined";
      project = "nixos";
    };
  };
}
