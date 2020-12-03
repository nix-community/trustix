{ stdenv
, lib
, buildGoModule
, fetchFromGitHub
}:
let
  pname = "trillian";
  version = "1.3.10";

in
buildGoModule {
  inherit pname version;
  vendorSha256 = "194zpzg6b7y046djs137rsiih2babzfzck1106c516xdpc2jj56m";

  src = fetchFromGitHub {
    owner = "google";
    repo = pname;
    rev = "v${version}";
    sha256 = "1cafrxg6p98n2y6n9cgyzrfg29wl0dsj65hqsqqcpijskcwwa04c";
  };

  # Remove tests that require networking
  postPatch = ''
    rm cmd/get_tree_public_key/main_test.go
  '';

  subPackages = [
    "cmd/trillian_log_server"
    "cmd/trillian_log_signer"
    "cmd/trillian_map_server"
    "cmd/createtree"
    "cmd/deletetree"
    "cmd/get_tree_public_key"
    "cmd/updatetree"
  ];

  meta = with lib; {
    homepage = "https://github.com/google/trillian";
    description = "A transparent, highly scalable and cryptographically verifiable data store.";
    license = [ licenses.asl20 ];
    maintainers = [ maintainers.adisbladis ];
  };
}
