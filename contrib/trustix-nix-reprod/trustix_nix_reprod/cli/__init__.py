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
