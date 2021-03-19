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

from trustix_nix_reprod.lib import DeferStack
from trustix_nix_reprod import (
    index_logs,
    index_eval,
    on_startup,
    on_shutdown,
)

import argparse
import asyncio

parser = argparse.ArgumentParser(
    description="Trustix Nix reproducibility dashboard CLI"
)

subparsers = parser.add_subparsers(dest="subcommand", required=True)

index_eval_parser = subparsers.add_parser("index_eval", help="Index evaluation")
index_eval_parser.add_argument(
    "--rev", default="nixos-unstable", help="Nixpkgs revision"
)

index_logs_parser = subparsers.add_parser(
    "index_logs", help="Index log build outputs (all known logs)"
)


def main():

    args = parser.parse_args()

    subcommand = args.subcommand

    async def _main():
        async with DeferStack() as defer:
            await on_startup()
            defer(on_shutdown)

            if subcommand == "index_eval":
                await index_eval(args.rev)
            elif subcommand == "index_logs":
                await index_logs()
            else:
                raise RuntimeError("Logic error")

    asyncio.run(_main())
