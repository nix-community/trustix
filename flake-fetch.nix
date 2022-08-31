# A minimal flake compatibility layer that only deals with fetching, and only deals with source types that we depend on
let
  inherit (builtins) fromJSON readFile mapAttrs filter attrNames hasAttr foldl';

  flake = fromJSON (readFile ./flake.lock);

  filterAttrs = fn: attrs:
    let
      filteredAttrs = filter (n: fn n attrs.${n}) (attrNames attrs);
    in
    foldl' (acc: attr: acc // { ${attr} = attrs.${attr}; }) { } filteredAttrs;

  fetchers = {

    github = locked: builtins.fetchTarball {
      sha256 = locked.narHash;
      url = "https://github.com/${locked.owner}/${locked.repo}/archive/${locked.rev}.tar.gz";
    };

  };

  fetchFlake = name: srcInfo:
    let
      inherit (srcInfo.locked) type;
      fetch = fetchers.${type};
    in
    fetch srcInfo.locked;

in
mapAttrs fetchFlake (filterAttrs (_: v: hasAttr "locked" v) flake.nodes)
