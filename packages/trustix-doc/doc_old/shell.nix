{ pkgs ? import <nixpkgs> { } }:
let
  tl = (pkgs.texlive.combine {
    inherit (pkgs.texlive) scheme-medium wrapfig ulem capt-of
      titlesec preprint enumitem paralist ctex environ svg
      beamer trimspaces zhnumber changepage framed pdfpages
      fvextra minted upquote ifplatform xstring;
  });

  pythonEnv = pkgs.python3.withPackages (ps: [
    ps.pygments
  ]);

in
pkgs.mkShell {
  buildInputs = [
    pythonEnv
    pkgs.pandoc
    pkgs.haskellPackages.pandoc-citeproc
    pkgs.haskellPackages.pandoc-crossref
    pkgs.plantuml
    tl
  ];
}
