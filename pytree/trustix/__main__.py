import pygit2 as git
import os.path

from . import lib

if __name__ == '__main__':

    repo_path = "./repo"

    r: git.Repository
    if os.path.exists(repo_path):
        r = lib.repo_open(repo_path)
        print("Opened")
    else:
        r = lib.repo_create(repo_path)
        print("Created")
