{ config, lib, pkgs, ... }:

let
  cfg = config.services.trustix-nix-r13y;
  inherit (lib) mkEnableOption mkOption types;

  # TODO: Get rid of hard coded port and use systemd sockets
  port = 8093;

  configFile = pkgs.writeText "trustix-nix-r13y.json" (builtins.toJSON {
    inherit (cfg) attrs channels lognames;
    log_poll_interval = cfg.logPollInterval;
  });

in
{
  options.services.trustix-nix-r13y = {
    enable = mkEnableOption (lib.mdDoc "the trustix-nix-r13y service");

    hostName = mkOption {
      type = types.str;
      description = lib.mdDoc ''
        Hostname to use for the service.
      '';
    };

    package = mkOption {
      type = types.package;
      default = pkgs.trustix-nix-r13y;
      defaultText = "pkgs.trustix-nix-r13y";
      description = "Which Trustix-Nix-r13y package to use.";
    };

    logPollInterval = mkOption {
      type = types.int;
      default = 900;
      description = lib.mdDoc ''
        How often to poll the trustix Daemon for published builds.
      '';
    };

    lognames = mkOption {
      type = types.attrsOf types.str;
      default = {};
      example = lib.literalExpression ''
        {
          "e0f263745e4e3ab07ab5275b00b44f594e0b6d2bd35892a8ebd10a7f86322eb7" = "trustix-demo";
        }
      '';
      description = ''
        Map log IDs to names in the UI.
      '';
    };

    attrs = mkOption {
      type = types.attrsOf (types.listOf types.str);
      default = [];
      example = lib.literalExpression ''
        {
          nixos-unstable = [ "hello" "jq" ];
        }
      '';
      description = ''
        Which attributes to show on landing page (grouped per channel).
      '';
    };

    channels = mkOption {
      default = {};
      description = lib.mdDoc ''
        Configuration for derivation imports.
      '';

      type = types.submodule {
        options = {

          hydra = mkOption {
            default = { };
            description = lib.mdDoc ''
              Hydra jobsets to import.
            '';
            type = types.attrsOf (types.submodule {
              options = {

                base_url = mkOption {
                  type = types.str;
                  description = lib.mdDoc ''
                    Hydra base URL.
                  '';
                };

                project = mkOption {
                  type = types.str;
                  description = lib.mdDoc ''
                    Hydra project.
                  '';
                };

                jobset = mkOption {
                  type = types.str;
                  description = lib.mdDoc ''
                    Hydra jobset.
                  '';
                };

                interval = mkOption {
                  type = types.int;
                  default = 3600;
                  description = lib.mdDoc ''
                    Poll interval (in seconds).
                  '';
                };

              };
            });
          };
        };
      };
    };
  };

  config = lib.mkIf cfg.enable {
    services.nginx.enable = true;

    # Static contents for the website
    services.nginx.virtualHosts.${cfg.hostName} = {
      root = pkgs.trustix-nix-r13y-web;
      extraConfig = ''
        index          index.html;
        try_files $uri /index.html;
      '';

      # Rewrite /api/foo requests to /foo before proxy passing
      locations."/api".extraConfig = ''
        rewrite /api/(.*) /$1  break;
        proxy_pass         http://localhost:${toString port};
        proxy_redirect     off;
        proxy_set_header   Host $host;
      '';
    };

    # Service
    systemd.services.trustix-nix-r13y = {
      description = "Trustix-Nix reproducibility tracker";
      wantedBy = [ "multi-user.target" ];
      # requires = [ "trustix-nix.socket" "trustix.socket" ];

      # binary-cache-proxy --address unix://${cfg.trustix-rpc} --listen ${cfg.listen}:${(toString cfg.port)} --privkey ${cfg.private-key}";
      serviceConfig = {
        Type = "simple";
        StateDirectory = "trustix-nix-r13y";
        ExecStart = "${lib.getBin cfg.package}/bin/trustix-nix-r13y serve --state /var/lib/trustix-nix-r13y --config ${configFile} --listen http://localhost:${toString port}";
        DynamicUser = true;
      };
    };

  };

}
