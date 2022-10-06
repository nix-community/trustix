let
  pubKey = "fG7JEPzIsr2mlSx5Xeh02BbJ4uzGpm5IE3aSGhS1UKo=";
in
{ config, pkgs, ... }: {
  imports = [
    ../packages/trustix/nixos
    ../packages/trustix-nix/nixos
  ];

  nixpkgs.pkgs = import ../pkgs.nix { };

  services.trustix = {
    enable = true;

    signers.snakeoil = {
      type = "ed25519";
      ed25519 = {
        # for testing this is fine, but in practice this should
        # be managed as a secret, and not put into the nix store
        private-key-path = pkgs.writeText "privkey" ''
          DyWQaOanQ64NU+k3dpp68/ABjFupTW941htRLRUCRdF8bskQ/MiyvaaVLHld6HTYFsni7MambkgTdpIaFLVQqg==
        '';
      };
    };

    publishers = [{
      signer = "snakeoil";
      protocol = "nix";
      publicKey.key = pubKey;
    }];
  };

  services.trustix-nix-build-hook = {
    enable = true;
    publisher = builtins.head config.services.trustix.publishers;
  };
}
