import subprocess
import json
import base64
import typing


B32_TRANS = str.maketrans(
    "0123456789abcdfghijklmnpqrsvwxyz", "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
)


def decode_nix_b32(encoded: str) -> bytes:
    return base64.b32decode(encoded.translate(B32_TRANS))


def parse_drv(drv_path) -> typing.Dict:
    p = subprocess.run(["nix", "show-derivation", drv_path], check=True, stdout=subprocess.PIPE)
    for d in json.loads(p.stdout).values():
        return d
