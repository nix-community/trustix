{ pkgs ? import ../../pkgs.nix { }
, lib ? pkgs.lib
}:

pkgs.stdenv.mkDerivation {
  pname = "trustix-doc";
  version = "dev";

  src = lib.cleanSource ./.;

  nativeBuildInputs = [
    pkgs.mdbook
  ];

  buildPhase = ''
    runHook preBuild

    ln -s ${lib.cleanSource ../trustix-proto} ../trustix-proto
    mdbook build

    runHook postBuild
  '';

  installPhase = ''
    mv book $out
  '';

}
