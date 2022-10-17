{ pkgs ? import ../../pkgs.nix { } }:

pkgs.mkShell {
  buildInputs = [
    pkgs.nodejs
  ];

  shellHook = ''
    export PATH=${builtins.toString ./.}/node_modules/.bin:$PATH
  '';
}
