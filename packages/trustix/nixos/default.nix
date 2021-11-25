{ config, lib, pkgs, ... }:

let
  cfg = config.services.trustix;

  configFile =
    let
      # toml doesn't have a "NoneType" so we must remove null attributes
      filterNull = attrs: lib.filterAttrsRecursive (n: v: v != null) attrs;
    in
    pkgs.writeText "trustix-config.json" (builtins.toJSON (
      filterNull (builtins.removeAttrs cfg [ "enable" "package" ])
    ));

  inherit (lib) mkOption types;

  deciderOpts =
    let
      mkDeciderModule = name: { ... }@options: mkOption {
        type = types.nullOr (types.submodule {
          inherit options;
        });
        description = "Configuration for the ${name} decision engine.";
        default = null;
      };
      mkDeciderModules = attrs: lib.mapAttrs (name: options: mkDeciderModule name options) attrs;

      deciderModules = {
        percentage = {
          minimum = mkOption {
            type = types.nullOr types.int;
            description = "Minimum agreement percentage.";
            default = null;
          };
        };

        javascript = {
          minimum = mkOption {
            type = types.nullOr types.lines;
            description = "JS script.";
            default = null;
          };
        };

        logid = {
          minimum = mkOption {
            type = types.nullOr types.str;
            description = "Configured log name to match.";
            default = null;
          };
        };
      };

    in
    {
      options = {

        engine = mkOption {
          type = types.enum (lib.mapAttrsToList (name: _: name) deciderModules);
          example = "percentage";
          description = "Which decision engine to use.";
        };

      } // mkDeciderModules deciderModules;
    };

  pubKeyOpts =
    {

      options = {

        type = mkOption {
          type = types.enum [ "ed25519" ];
          example = "ed25519";
          default = "ed25519";
          description = "Key type.";
        };

        key = mkOption {
          type = types.str;
          example = "2uy8gNIOYEewTiV7iB7cUxBGpXxQtdlFepFoRvJTCJo=";
          default = "Base64 encoded public key";
          description = "Key data.";
        };
      };

    };

  publisherOpts =
    {
      options = {

        protocol = mkOption {
          type = types.str;
          example = "nix";
          description = "Subprotocol id/name.";
        };

        publicKey = mkOption {
          type = types.submodule pubKeyOpts;
          description = "Decision making engine configurations.";
        };

        meta = mkOption {
          type = types.attrsOf types.str;
          default = { };
          description = ''
            Arbitrary metadata to set for a log.
          '';
        };

        signer = mkOption {
          type = types.str;
          example = "snakeoil";
          description = "Signer name.";
        };

      };
    };

  signerOpts =
    {
      options = {

        type = mkOption {
          type = types.enum [ "ed25519" ];
          example = "ed25519";
          default = "ed25519";
          description = "Signing backend.";
        };

        ed25519 = {
          private-key-path = mkOption {
            type = types.path;
            description = "Path to private key.";
          };
        };

      };
    };

  subscriberOpts =
    {
      options = {

        protocol = mkOption {
          type = types.str;
          example = "nix";
          description = "Subprotocol id/name.";
        };

        publicKey = mkOption {
          type = types.submodule pubKeyOpts;
          description = "Decision making engine configurations.";
        };

        meta = mkOption {
          type = types.attrsOf types.str;
          default = { };
          description = ''
            Arbitrary metadata to set for a log.
          '';
        };

        syncMode = mkOption {
          type = types.enum [ "light" ];
          default = "light";
          description = "Verification and data sync mode";
        };

      };
    };


in
{

  options.services.trustix = {

    enable = lib.mkEnableOption "trustix";

    package = mkOption {
      type = types.package;
      default = pkgs.trustix;
      defaultText = "pkgs.trustix";
      description = "Which Trustix derivation to use.";
    };

    deciders = mkOption {
      type = types.attrsOf (types.submodule deciderOpts);
      default = { };
      description = "Decision making engine configurations (scoped per subprotocol).";
    };

    signers = mkOption {
      type = types.attrsOf (types.submodule signerOpts);
      default = { };
      description = "Log signers for published logs.";
    };

    publishers = mkOption {
      type = types.listOf (types.submodule publisherOpts);
      default = [ ];
      description = "Publisher configurations.";
    };

    subscribers = mkOption {
      type = types.listOf (types.submodule subscriberOpts);
      default = [ ];
      description = "Subscriber configurations.";
    };

    storage = {
      type = mkOption {
        type = types.enum [ "native" ];
        default = "native";
        internal = true;
        description = "Storage engine.";
      };
    };

    remotes = mkOption {
      type = types.listOf types.str;
      default = [ ];
      description = "List of remotes to connect to.";
    };

  };

  config = lib.mkIf cfg.enable {

    users.users.trustix = {
      isSystemUser = true;
      description = "The user that the trustix daemon runs as.";
      group = "trustix";
    };

    users.groups.trustix = { };

    systemd.sockets.trustix = {
      description = "Socket for the Trustix daemon";
      wantedBy = [ "sockets.target" ];
      listenStreams = [ "/run/trustix-daemon.socket" ];
    };

    systemd.services.trustix = {
      description = "Trustix daemon";
      wantedBy = [ "multi-user.target" ];
      requires = [ "trustix.socket" ];

      serviceConfig = {
        Type = "simple";
        User = "trustix";
        Group = "trustix";
        ExecStart = "${lib.getBin cfg.package}/bin/trustix daemon --state . --config ${configFile}";
        StateDirectory = "trustix";
        WorkingDirectory = "%S/trustix";
        DynamicUser = true;
      };
    };

  };

}
