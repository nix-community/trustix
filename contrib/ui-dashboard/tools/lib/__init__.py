import subprocess
import os.path
import typing
import time
import os


_SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))

TOOLS_DIR = os.path.dirname(_SCRIPT_DIR)
ROOT_DIR = os.path.dirname(TOOLS_DIR)
STATE_DIR = os.path.join(ROOT_DIR, "state")

PSQL_DATA_DIR = os.path.join(STATE_DIR, "psql-data")
PSQL_SOCKETS_DIR = os.path.join(os.environ["TMPDIR"], "nix-trustix-dash-psql-sockets")
PSQL_DB_NAME = "nix-trustix-dash"


def ensure_dir(path: str):
    if not os.path.exists(PSQL_SOCKETS_DIR):
        os.mkdir(PSQL_SOCKETS_DIR)


def wait_for_psql():
    """Wait for postgresql to be up and running"""
    socket = os.path.join(PSQL_SOCKETS_DIR, ".s.PGSQL.5432")
    while not os.path.exists(socket):
        time.sleep(0.1)


def wait_for_db():
    """Wait for postgresql to be up and running and the database to be created"""
    wait_for_psql()

    while True:
        p = subprocess.run(
            ["psql", "-h", PSQL_SOCKETS_DIR, PSQL_DB_NAME, "-c", "\\q"],
            stderr=subprocess.PIPE,
        )
        if p.returncode == 0:
            break
        time.sleep(0.1)


def get_fmt_files() -> typing.List[str]:
    """Return a list of files/directories to format using black"""
    ret: typing.List[str] = [
        os.path.join(TOOLS_DIR, f)
        for f in os.listdir(TOOLS_DIR)
        if not f.startswith(".") and not f.startswith("#") and not os.path.isdir(f)
    ]
    ret.append(ROOT_DIR)
    return ret


def exec_cmd(cmdline: typing.List[str]):
    os.execvp(cmdline[0], cmdline)
