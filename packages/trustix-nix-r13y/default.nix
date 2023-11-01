{ buildGoApplication
, lib
, gitignoreSource
, makeWrapper
, nix-eval-jobs
, nix
, diffoscope
, git
, stdenv
}:

let
  runtimeDeps = [
    nix
    nix-eval-jobs
    (diffoscope.override {
      enableBloat = ! stdenv.isDarwin;
    })
    git
  ];

  inherit (lib) hasSuffix;

in
buildGoApplication {
  pname = "trustix-nix-r13y";
  version = "dev";
  pwd = ./.;
  modules = ./gomod2nix.toml;
  src = lib.cleanSourceWith {
    filter = name: type: (
      ! hasSuffix "tests" name
      && ! hasSuffix "nixos" name
    );
    src = gitignoreSource ./.;
  };

  nativeBuildInputs = [
    makeWrapper
  ];

  postInstall = ''
    wrapProgram "$out/bin/trustix-nix-r13y" --prefix PATH : "${lib.makeBinPath runtimeDeps}"
  '';
}
