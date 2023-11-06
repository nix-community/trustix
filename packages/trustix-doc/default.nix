{ pkgs
, lib
}:

pkgs.stdenv.mkDerivation {
  pname = "trustix-doc";
  version = "dev";

  src = ../..;

  nativeBuildInputs = [
    pkgs.mdbook
  ];

  buildPhase = ''
    runHook preBuild

    cd packages/trustix-doc
    mdbook build

    runHook postBuild
  '';

  installPhase = ''
    mv book $out
  '';

}
