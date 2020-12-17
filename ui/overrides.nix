{ pkgs }: self: super: {

  celery = super.celery.overridePythonAttrs (old: {
    propagatedBuildInputs = old.propagatedBuildInputs ++ [ self.setuptools ];
  });

  jsonslicer = super.jsonslicer.overridePythonAttrs (old: {
    nativeBuildInputs = old.nativeBuildInputs ++ [ pkgs.pkgconfig ];
    buildInputs = old.buildInputs ++ [ pkgs.yajl ];
  });

  packaging = super.packaging.overridePythonAttrs (old: {
    buildInputs = [ self.flit-core ];
  });

  pynix = super.pynix.overridePythonAttrs (old: {
    buildInputs = [ self.poetry ];
  });

  daphne = super.daphne.overridePythonAttrs (old: {
    postPatch = ''
      substituteInPlace setup.py --replace 'setup_requires=["pytest-runner"],' ""
    '';
  });

}
