import pygit2 as git
import os.path
import typing
import time


def repo_create(path: str) -> git.Repository:
    """Create bare repository with partial clone enabled"""
    r = git.init_repository(path, bare=True)

    c = r.config
    # Note that dict mutation is overloaded and implies saving
    c["uploadpack.allowanysha1inwant"] = True
    c["uploadpack.allowfilter"] = True

    return r


def repo_open(path: str) -> git.Repository:
    """Open repo and assert partial clone is enabled"""
    r = git.Repository(path)

    c = r.config
    if c["uploadpack.allowanysha1inwant"] != "true" or c["uploadpack.allowfilter"] != "true":
        raise ValueError("Invalid repo configuration, not configured for partial cloning")

    return r


class Repository:

    _repo: git.Repository
    _tree: typing.Optional[typing.Any]
    _commit: typing.Any
    _name: str
    _email: str

    def __init__(self, repo_path: str, name: str, email: str):

        self._name = name
        self._email = email

        if os.path.exists(repo_path):
            self._repo = repo_open(repo_path)
            self._commit = self._repo.head.resolve().target
            self._tree = self._repo.get(self._commit).tree.id

        else:
            self._commit = None
            self._repo = repo_create(repo_path)
            self._tree = self._repo.TreeBuilder().write()
            self.write_commit(message="Init log")

    def write_commit(self, message: typing.Optional[str] = ""):
        now = int(time.time())

        parents: typing.List[git.Oid]
        if self._commit:
            parents = [ self._commit ]
        else:
            parents = []

        sig = git.Signature(self._name, self._email, time=now)

        self._commit = self._repo.create_commit(
            "HEAD",
            sig,
            sig,
            message,
            self._tree,
            parents,
        )
