{ config, pkgs, ... }:

{
  services.nginx.virtualHosts."trustix.dev" = {
    root = pkgs.trustix-doc;
  };
}
