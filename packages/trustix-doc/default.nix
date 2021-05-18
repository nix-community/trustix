{ pkgs ? import ../../nix, lib ? pkgs.lib }:

let
  rootShell = import ../../shell.nix;

in pkgs.stdenv.mkDerivation {
  pname = "trustix-doc";
  version = "dev";

  src = ./.;

  nativeBuildInputs = [
    pkgs.mdbook
  ];

  preBuild = ''
    ln -s ${lib.cleanSource ../trustix-proto} ../trustix-proto
  '';

  installPhase = ''
    mv book $out
  '';

}
