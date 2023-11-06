{ config, options, lib, pkgs, ... }:

let
  cfg = config.services.trustix-nix-build-hook;

  hook-script = pkgs.writeShellScript "trustix-hook"
    ''
      set -euo pipefail
      export TRUSTIX_TOKEN="${cfg.token-path}"
      LOG_ID=${
        if builtins.isString cfg.publisher
        then cfg.publisher
        else builtins.concatStringsSep " " [
          "$("
          "${lib.getBin pkgs.trustix}/bin/trustix"
          "print-log-id"
          "--protocol" cfg.publisher.protocol
          "--pubkey" cfg.publisher.publicKey.key
          ")"
        ]
      }
      ${lib.getBin pkgs.trustix-nix}/bin/trustix-nix --log-id $LOG_ID post-build-hook --address unix://${cfg.trustix-rpc}
    '';

  inherit (lib) mkOption types literalExpression;
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

    logID = mkOption {
      type = types.str;
      description = ''
        DEPRECATED, use `publisher`
        Which local Trustix log to submit build results to.
      '';
    };

    publisher = mkOption {
      type = types.either types.str (options.services.trustix.publishers.type.nestedTypes.elemType or types.unspecified);
      description = "Which local Trustix log to submit build results to.";
      example = literalExpression "builtins.head config.services.trustix.publishers";
    };

    trustix-rpc = mkOption {
      type = types.path;
      default = "/run/trustix-daemon.socket";
      description = "Which Trustix socket to connect to.";
    };

    token-path = mkOption {
      type = types.path;
      default = "/var/lib/trustix/trustix.token";
      description = "Path to write token.";
    };
  };

  config = lib.mkIf cfg.enable {

    nix.extraOptions = ''
      post-build-hook = ${hook-script}
    '';

    services.trustix-nix-build-hook.publisher = lib.mkDerivedConfig
      options.services.trustix-nix-build-hook.logID
      (x: x);
  };

}
