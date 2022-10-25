# Find all passthru.tests and return them as a list of full attribute paths
builtins.toJSON (
  let
    p = import ../. { };
  in
  builtins.foldl'
    (acc: v:
    let
      tests = builtins.attrNames (p.${v}.passthru.tests or { });
    in
    acc ++ map (t: "${v}.passthru.tests.${t}") tests) [ ]
    (builtins.attrNames p)
)
