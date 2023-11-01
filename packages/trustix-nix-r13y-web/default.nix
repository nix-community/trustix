{ npmlock2nix, lib, gitignoreSource, nodejs }:

npmlock2nix.v2.build {
  src = gitignoreSource ./.;
  installPhase = "cp -r dist $out";
  buildCommands = [ "env HOME=$(mktemp -d) npm run build" ];
  inherit nodejs;
}
