{ npmlock2nix, lib, gitignoreSource, nodejs }:

npmlock2nix.v2.build {
  src = gitignoreSource ./.;
  installPhase = "cp -r dist $out";
  buildCommands = [ "npm run build" ];
  inherit nodejs;
}
