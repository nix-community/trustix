import typing
import base64
import ast


B32_TRANS = str.maketrans(
    "0123456789abcdfghijklmnpqrsvwxyz", "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
)


def decode_nix_b32(encoded: str) -> bytes:
    return base64.b32decode(encoded.translate(B32_TRANS))


def get_drv_refs(drv_path: str) -> typing.List[str]:
    """Return all build-time references for a derivation"""

    refs: typing.List[str] = []

    parsed = ast.parse(open(drv_path).read())
    for x in parsed.body[0].value.args[1].elts:  # type: ignore
        for y in x.elts:
            if isinstance(y, ast.Constant):
                refs.append(y.value)
    return refs
