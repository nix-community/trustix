{ stdenv, buildNodeModules, lib, nodejs, npmHooks }:

stdenv.mkDerivation {
  pname = "trustix-nix-r13y-web";
  version = "0.1.0";

  src = ./.;

  nativeBuildInputs = [
    buildNodeModules.hooks.npmConfigHook
    nodejs
    npmHooks.npmInstallHook
  ];

  nodeModules = buildNodeModules.fetchNodeModules {
    packageRoot = ./.;
  };
}
