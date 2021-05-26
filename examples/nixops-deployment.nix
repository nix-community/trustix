{

  network.description = "Example Trustix NixOps deployment";

  publisher_server =
    { resources, ... }:
    {

      # Add the root-level NixOS module
      imports = [
        /path/to/trustix/git/checkout/nixos
      ];

      services.trustix = {
        enable = true;

        signer.snakeoil = {
          type = "ed25519";
          ed25519 = {
            # Keys from deployment.keys are stored under /run/ on a temporary filesystem and will not persist across a reboot.
            private-key-path = "/run/keys/trustix-priv";
          };
        };

        publishers = [
          {
            signer = "snakeoil";
            protocol = "nix";
            key = {
              type = "ed25519";
              # Contents of the generated trustix-pub
              pub = "2uy8gNIOYEewTiV7iB7cUxBGpXxQtdlFepFoRvJTCJo=";
            };
          }
        ];
      };

      # Push local builds via the post-build hook
      services.trustix-nix-build-hook.enable = true;

      deployment = {
        # Replace with your desired deployment target
        # For example ec2/gcp
        #
        # targetHost = "10.0.0.1";

        # Upload keys to the remote
        keys = {
          trustix-priv.keyFile = "./trustix-priv";
        };
      };
    };

  subscriber_server =
    { resources, ... }:
    {

      services.trustix = {
        enable = true;

        subscribers = [
          {
            protocol = "nix";
            key = {
              type = "ed25519";
              # Contents of the generated trustix-pub
              pub = "2uy8gNIOYEewTiV7iB7cUxBGpXxQtdlFepFoRvJTCJo=";
            };
          }
        ];

        # A remote can expose many logs and they are not neccesarily created by the remote in question
        remotes = [
          "grpc+http://10.0.0.1"
        ];

      };

      deployment = {
        # Replace with your desired deployment target
        # For example ec2/gcp
        #
        # targetHost = "10.0.0.2";
      };

    };

}
