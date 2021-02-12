#!/usr/bin/env python
from pprint import pprint
import subprocess
import tempfile
import os.path
import typing
import shutil
import black
import glob
import sys

# gRPC/protobuf python bindings lack support for something equivalent to go_package to set a custom package name
# Hence we need to rewrite the imports so that we don't get package names like "api", "proto" & "schema"


PACKAGES: typing.List[str] = [
    "proto",
    "schema",
    "api",
]
NAME_PREFIX = "trustix"


def dirname_recurse(filepath: str, depth: int) -> str:
    """Find a parent directory by relative depth"""
    for _ in range(depth):
        filepath = os.path.dirname(filepath)
    return filepath


if __name__ == "__main__":
    script_path = os.path.abspath(__file__)
    root_dir = dirname_recurse(script_path, 4)

    with tempfile.TemporaryDirectory() as tmp_dir:
        for package in PACKAGES:

            # Generate python files
            cmdline: typing.List[str] = [
                sys.executable,
                "-m",
                "grpc_tools.protoc",
                "-I",
                root_dir,
                "-I",
                os.path.join(root_dir, package),
                f"--python_out={tmp_dir}",
                f"--grpc_python_out={tmp_dir}",
            ] + glob.glob(os.path.join(root_dir, package, "*.proto"))
            subprocess.run(cmdline, check=True)

            if not os.path.exists(os.path.join(tmp_dir, package, "__init__.py")):
                with open(os.path.join(tmp_dir, package, "__init__.py"), "w") as f:
                    f.write("")

            # Iterate over python files and rewrite imports to be prefix by trustix_
            # Technically an AST aware replacement is "better" but a simple string replacement is good enough
            for py_file in glob.glob(os.path.join(tmp_dir, package, "*.py")):
                with open(py_file) as f:
                    contents = f.read()

                for import_pkg in PACKAGES:
                    contents = contents.replace(
                        f"from {import_pkg} import",
                        f"from {NAME_PREFIX}_{import_pkg} import",
                    )

                # Pre-format
                try:
                    contents = black.format_file_contents(
                        contents, fast=False, mode=black.FileMode()
                    )
                except black.NothingChanged:
                    pass

                with open(py_file, "w") as f:
                    f.write(contents)

            # Move the package into the correct directory in the source tree
            rewritten_package = f"{NAME_PREFIX}_{package}"
            package_path = os.path.join("..", rewritten_package)

            if os.path.exists(package_path):
                shutil.rmtree(package_path)

            shutil.move(os.path.join(tmp_dir, package), package_path)
