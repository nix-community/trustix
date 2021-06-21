{ config, lib, pkgs, ... }:

let
  cfg = config.services.trustix-nix-build-hook;

  hook-script = pkgs.writeScript "trustix-hook"
    ''
      ${lib.getBin pkgs.trustix-nix}/bin/trustix-nix post-build-hook --address ${cfg.trustix-rpc}
    '';

  inherit (lib) mkOption types;
in
{

  options.services.trustix-nix-build-hook = {

    enable = lib.mkEnableOption "Trustix Nix post-build hook";

    package = mkOption {
      type = types.package;
      default = pkgs.trustix-nix;
      defaultText = "pkgs.trustix-nix";
      description = "Which Trustix-Nix derivation to use.";
    };

    trustix-rpc = mkOption {
      type = types.path;
      default = "/run/trustix-daemon.socket";
      description = "Which Trustix socket to connect to.";
    };

  };

  config = lib.mkIf cfg.enable {
    nix.extraOptions = ''
      post-build-hook = ${hook-script}
    '';
  };

}
