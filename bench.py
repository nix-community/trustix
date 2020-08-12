#!/usr/bin/env python
from concurrent.futures import ThreadPoolExecutor
import subprocess
import textwrap
import time
import os


CMD = "./dev/test"
WORKERS = 10
JOBS = 100000
JOBS_PER_WORKER = int(JOBS / WORKERS)


def run():
    devnull = open(os.devnull, 'w')
    processes = []

    for _ in range(JOBS_PER_WORKER):
        p = subprocess.Popen(CMD, stdout=devnull, stderr=devnull)
        processes.append(p)

    for p in processes:
        p.wait()


if __name__ == '__main__':
    start = time.time()

    with ThreadPoolExecutor(max_workers=WORKERS) as e:
        e.submit(run)

    end = time.time() - start

    print(textwrap.dedent(f"""
    Jobs: {JOBS}
    Time: {end}
    Jobs/second: {JOBS/end}
    """).strip())
