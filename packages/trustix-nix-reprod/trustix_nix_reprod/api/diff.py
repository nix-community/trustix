# Trustix
# Copyright (C) 2021 Tweag IO

# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

from concurrent.futures import ThreadPoolExecutor
from typing import Dict
import subprocess
import requests
import tempfile
import asyncio
import codecs
import json
import os

from trustix_nix_reprod.cache import cached
from trustix_python.api import api_pb2
from trustix_nix_reprod.conf import settings
from trustix_nix_reprod.api.models import DiffResponse
from trustix_nix_reprod.proto import get_rpcapi


# Uvloop has a nasty bug https://github.com/MagicStack/uvloop/issues/317
# To work around this we run the fetching/unpacking in a separate blocking thread
def _fetch_unpack_nar(url, location):
    loc_dir = os.path.dirname(location)

    try:
        os.mkdir(loc_dir)
    except FileExistsError:
        pass

    with requests.get(url, stream=True) as r:
        r.raise_for_status()
        p = subprocess.Popen(
            ["nix-store", "--restore", loc_dir],
            stdin=subprocess.PIPE,
        )
        for chunk in r.iter_content(chunk_size=512):
            p.stdin.write(chunk)
        p.stdin.close()
        p.wait(timeout=0.5)

    # Ensure correct mtime
    for subl in (
        (os.path.join(dirpath, f) for f in (dirnames + filenames))
        for (dirpath, dirnames, filenames) in os.walk(location)
    ):
        for path in subl:
            os.utime(path, (1, 1))
    os.utime(location, (1, 1))


def _process_narinfo(narinfo: Dict, tmpdir, outbase) -> str:
    nar_hash = narinfo["narHash"].split(":")[-1]
    store_base = narinfo["path"].split("/")[-1]

    store_prefix = store_base.split("-")[0]

    unpack_dir = os.path.join(tmpdir, store_base, outbase)
    nar_url = "/".join((settings.binary_cache_proxy, "nar", store_prefix, nar_hash))

    _fetch_unpack_nar(nar_url, unpack_dir)

    return unpack_dir


def _diff(narinfo1: Dict, narinfo2: Dict) -> Dict:
    with tempfile.TemporaryDirectory(prefix="trustix-ui-dash-diff") as tmpdir:
        with ThreadPoolExecutor(max_workers=2) as e:
            dir_a_fut = e.submit(_process_narinfo, narinfo1, tmpdir, "A")
            dir_b_fut = e.submit(_process_narinfo, narinfo2, tmpdir, "B")
            dir_a = dir_a_fut.result()
            dir_b = dir_b_fut.result()

        dir_a_rel = os.path.join(os.path.basename(os.path.dirname(dir_a)), "A")
        dir_b_rel = os.path.join(os.path.basename(os.path.dirname(dir_b)), "B")

        proc = subprocess.run(
            ["diffoscope", "--json", "-", dir_a_rel, dir_b_rel],
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
            cwd=tmpdir,
        )

    # Diffoscope returns non-zero on paths that have a diff
    # Instead use stderr as a heurestic if the call went well or not
    if proc.stderr:
        raise ValueError(proc.stderr)

    return json.loads(proc.stdout)


@cached(model=DiffResponse, ttl=settings.cache_ttl.diff)
async def diff(output_hash_1_hex: str, output_hash_2_hex: str) -> DiffResponse:
    output_hash_1 = codecs.decode(output_hash_1_hex, "hex")  # type: ignore
    output_hash_2 = codecs.decode(output_hash_2_hex, "hex")  # type: ignore

    rpc_client = get_rpcapi()

    narinfo1, narinfo2 = [
        json.loads(resp.Value)
        for resp in (
            await asyncio.gather(
                rpc_client.GetValue(api_pb2.ValueRequest(Digest=output_hash_1)),  # type: ignore
                rpc_client.GetValue(api_pb2.ValueRequest(Digest=output_hash_2)),  # type: ignore
            )
        )
    ]

    diffoscope = await asyncio.get_running_loop().run_in_executor(
        None, _diff, narinfo1, narinfo2
    )

    return DiffResponse(
        narinfo={
            output_hash_1_hex: narinfo1,
            output_hash_2_hex: narinfo2,
        },
        diffoscope=diffoscope,
    )
