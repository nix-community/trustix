{ config, lib, pkgs, ... }:

let
  cfg = config.services.trustix;

  configFile =
    let
      configJSON = pkgs.writeText "trustix-config.json" (builtins.toJSON (
        builtins.removeAttrs cfg [ "enable" "package" ]
      ));
    in
    pkgs.runCommand "trustix-config.toml"
      {
        nativeBuildInputs = [ pkgs.remarshal ];
        preferLocalBuild = true;
      } ''
      remarshal -i ${configJSON} --if json -o $out --of toml
    '';

  inherit (lib) mkOption types;

  deciderOpts =
    let
      mkDeciderModule = name: { ... }@options: mkOption {
        type = types.submodule {
          inherit options;
        };
        description = "Configuration for the ${name} decision engine.";
      };
    in
    {
      options = {

        engine = mkOption {
          type = types.enum [ "percentage" "lua" "logid" ];
          example = "percentage";
          description = "Which decision engine to use.";
        };

        percentage = mkDeciderModule "percentage" {
          minimum = mkOption {
            type = types.nullOr types.int;
            description = "Minimum agreement percentage.";
            default = null;
          };
        };

        javascript = mkDeciderModule "javascript" {
          minimum = mkOption {
            type = types.nullOr types.lines;
            description = "JS script.";
            default = null;
          };
        };

        logid = mkDeciderModule "logid" {
          minimum = mkOption {
            type = types.nullOr types.str;
            description = "Configured log name to match.";
            default = null;
          };
        };

      };
    };

  pubKeyOpts =
    {

      type = mkOption {
        type = types.enum [ "ed25519" ];
        example = "ed25519";
        default = "ed25519";
        description = "Key type.";
      };

      pub = mkOption {
        type = types.str;
        example = "2uy8gNIOYEewTiV7iB7cUxBGpXxQtdlFepFoRvJTCJo=";
        default = "Base64 encoded public key";
        description = "Key data.";
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
          type = types.AttrsOf types.str;
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
          type = types.AttrsOf types.str;
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
        ExecStart = "${lib.getBin cfg.package}/bin/trustix daemon --state . --config ${configFile}";
        StateDirectory = "trustix";
        WorkingDirectory = "%S/trustix";
        DynamicUser = true;
      };
    };

  };

}
