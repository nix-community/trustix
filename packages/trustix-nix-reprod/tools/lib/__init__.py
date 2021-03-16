from watchdog.events import FileSystemEventHandler  # type: ignore
from watchdog.observers import Observer  # type: ignore
import subprocess
import threading
import traceback
import os.path
import typing
import time
import sys
import os


_SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))

TOOLS_DIR = os.path.dirname(_SCRIPT_DIR)
ROOT_DIR = os.path.dirname(TOOLS_DIR)
STATE_DIR = os.environ["NIX_REPROD_STATE_DIR"]

PSQL_DATA_DIR = os.path.join(STATE_DIR, "psql-data")
PSQL_SOCKETS_DIR = os.path.join(os.environ["TMPDIR"], "nix-trustix-reprod-psql-sockets")
PSQL_DB_NAME = "nix-trustix-reprod"
PSQL_DB_URI = f"postgres:///{PSQL_DB_NAME}?host={PSQL_SOCKETS_DIR}"


_icons = [
    "🤪",
    "👻",
    "😁",
    "🥕",
    "🌮",
    "🍆",
    "🥦",
    "🔥",
    "🙃",
    "🥳",
    "🥸",
    "🧙",
    "🚀",
]


def ensure_dir(path: str):
    if not os.path.exists(PSQL_SOCKETS_DIR):
        os.mkdir(PSQL_SOCKETS_DIR)


def wait_for_psql():
    """Wait for postgresql to be up and running"""
    socket = os.path.join(PSQL_SOCKETS_DIR, ".s.PGSQL.5432")
    while not os.path.exists(socket):
        time.sleep(0.1)


def db_exists() -> bool:
    p = subprocess.run(
        ["psql", "-h", PSQL_SOCKETS_DIR, PSQL_DB_NAME, "-c", "\\q"],
        stderr=subprocess.PIPE,
    )
    return p.returncode == 0


def wait_for_db():
    """Wait for postgresql to be up and running and the database to be created"""
    wait_for_psql()

    while True:
        if db_exists():
            break
        time.sleep(0.1)


def get_watch_files() -> typing.List[str]:
    return [
        os.path.join(TOOLS_DIR, f)
        for f in os.listdir(TOOLS_DIR)
        if not f.startswith(".") and not f.startswith("#") and not os.path.isdir(f)
    ] + [
        os.path.join(ROOT_DIR, f)
        for f in (
            "pynix",
            os.path.join("tools", "lib"),
        )
    ]


def exec_cmd(cmdline: typing.List[str]):
    os.execvp(cmdline[0], cmdline)


def run_cmd(cmdline: typing.List[str]) -> int:
    p = subprocess.run(cmdline)
    return p.returncode


def watch_recursive(
    files: typing.List[str], handler: typing.Callable, delay: float = 0.5
):

    evt = threading.Event()

    def handler_loop():
        while True:
            evt.wait()
            time.sleep(delay)
            evt.clear()
            try:
                handler()
            except Exception:
                traceback.print_exc(file=sys.stderr)

    class WatchHandler(FileSystemEventHandler):
        def on_any_event(self, event):
            evt.set()

    observer = Observer()

    handler_thread = threading.Thread(target=handler_loop, daemon=True)
    handler_thread.start()

    for f in files:
        observer.schedule(WatchHandler(), f, recursive=os.path.isdir(f))

    observer.start()
    try:
        while True:
            time.sleep(0.1)
    finally:
        observer.stop()
        observer.join()


def _djb2_hash(s: str) -> int:
    h: int = 5381
    for x in s:
        h = ((h << 5) + h) + ord(x)
    return h & 0xFFFFFFFF


def icon(s: str) -> str:
    """Return a deterministic icon based on an input string"""
    return _icons[_djb2_hash(s) % len(_icons)]
