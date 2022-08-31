{ pkgs }: self: super: {

  aerich = super.aerich.overridePythonAttrs (old: {
    nativeBuildInputs = old.nativeBuildInputs ++ [
      self.poetry
    ];
  });

  # python-magic = super.python-magic.overridePythonAttrs (old: {
  #   inherit (pkgs.python3Packages.python_magic) patches;
  # });

  libarchive-c = super.libarchive-c.overridePythonAttrs (old: {
    postPatch = ''
      substituteInPlace libarchive/ffi.py --replace \
        "find_library('archive')" "'${pkgs.libarchive.lib}/lib/libarchive${pkgs.stdenv.hostPlatform.extensions.sharedLibrary}'"
    '';
  });

}
