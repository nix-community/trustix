{ pkgs
, lib
, gitignoreSource
}:

pkgs.stdenv.mkDerivation {
  pname = "trustix-doc";
  version = "dev";

  src = lib.cleanSourceWith {
    filter = name: type: ! lib.hasSuffix "tests" name;
    src = gitignoreSource ./.;
  };

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
