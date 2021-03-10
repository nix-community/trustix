import argparse
import asyncio

parser = argparse.ArgumentParser(
    description="Trustix Nix reproducibility dashboard CLI"
)

subparsers = parser.add_subparsers(dest="subcommand", required=True)

index_eval_parser = subparsers.add_parser("index_eval", help="Index evaluation")
index_eval_parser.add_argument("--rev", default="nixos-unstable", help="Nixpkgs revision")

index_log_parser = subparsers.add_parser("index_log", help="Index log build outputs")
index_log_parser.add_argument("log_name", help="Log name")

index_logs_parser = subparsers.add_parser("index_logs", help="Index log build outputs (all known logs)")


def main():

    args = parser.parse_args()

    subcommand = args.subcommand

    async def _main():
        if subcommand == "index_eval":
            pass
        elif subcommand == "index_log":
            pass
        elif subcommand == "index_logs":
            pass
        else:
            raise RuntimeError("Logic error")

    asyncio.run(_main())
