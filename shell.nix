let
  pkgs = import ./nix;

  STATE_DIR = "${builtins.toString ./.}/state";
  TRUSTIX_RPC = "unix://${STATE_DIR}/trustix.sock";

in pkgs.mkShell {

  buildInputs = [
    pkgs.nixpkgs-fmt
    pkgs.hivemind
    pkgs.niv
  ];

  inherit STATE_DIR TRUSTIX_RPC;

}
