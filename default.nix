{ pkgs ? import <nixpkgs> {} }:

let
  gitignoreSrc = pkgs.fetchFromGitHub {
    owner = "hercules-ci";
    repo = "gitignore";
    rev = "6c4ab20841d2a20cf69d52c8e848c4d6b0aa73fe";
    sha256 = "0nbwg01z0girs8c5zxg5zqivhny064rafzf47pygsmxlag0jiliq";
  };
  inherit (import gitignoreSrc { inherit (pkgs) lib; }) gitignoreSource;

  inherit (pkgs) stdenv buildGoPackage darwin;

in null # buildGoPackage {
#   pname = "trustix";
#   version = "git";

#   goDeps = ./deps.nix;

#   goPackagePath = "github.com/adisbladis/trustix";
#   # Fix for usb-related segmentation faults on darwin
#   propagatedBuildInputs =
#     stdenv.lib.optionals stdenv.isDarwin [ darwin.libobjc darwin.IOKit ];

#   nativeBuildInputs = [ pkgs.go-ethereum pkgs.solc ];

#   # Fixes Cgo related build failures (see https://github.com/NixOS/nixpkgs/issues/25959 )
#   hardeningDisable = [ "fortify" ];

#   src = gitignoreSource ./.;

#   preConfigure = ''
#     make contract
#   '';

# }
