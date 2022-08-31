{ system ? builtins.currentSystem }:
let
  flake = builtins.getFlake (builtins.toString ./.);
in
flake.packages.${system}
