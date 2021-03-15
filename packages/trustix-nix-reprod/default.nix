{ poetry2nix
, nix-nar-unpack
, hydra-eval-jobs
, diffoscope
, pkgs
}:

poetry2nix.mkPoetryApplication {
  projectDir = ./.;

  propagatedBuildInputs = [ hydra-eval-jobs nix-nar-unpack diffoscope ];

  # Don't propagate anything, hydra-eval-jobs is already wrapped in $PATH
  postFixup = "rm $out/nix-support/propagated-build-inputs";

  overrides = poetry2nix.overrides.withDefaults (
    import ./overrides.nix { inherit pkgs; }
  );
}
