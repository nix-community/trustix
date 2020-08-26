import pygit2 as git  # type: ignore
import os.path
import hashlib
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


def shard(input: str) -> typing.Tuple[str, ...]:
    """
    Decide where an input to the tree ends up
    """
    depth = 6

    return tuple(
        input[2 * i : 2 * i + 2]
        for i in range(depth)
    ) + (input[2 * depth:],)


def auto_insert(repo, treebuilder, path, content):
    tree_oid = treebuilder.write()
    tree = repo.get(tree_oid)

    if len(path) == 1:  # Reached leaf, insert blob
        blob = repo.create_blob(content)
        treebuilder.insert(path[0], blob, git.GIT_FILEMODE_BLOB)

        names = set(e.name for e in tree)
        names.add(path[0])

        m = hashlib.sha256()

        for n in sorted(names):
            if n == path[0]:
                data = content
            else:
                data = tree[n].read_raw()
            m.update(data)

        treebuilder.insert("hash", repo.create_blob(m.digest()), git.GIT_FILEMODE_BLOB)
        return treebuilder.write()

    subtree_name, sub_path = path[0], path[1:]
    try:
        sub_treebuilder = repo.TreeBuilder(tree[subtree_name].oid)
    except KeyError:
        sub_treebuilder = repo.TreeBuilder()

    # Update or create nested subtree
    subtree_oid = auto_insert(repo, sub_treebuilder, sub_path, content)
    treebuilder.insert(subtree_name, subtree_oid, git.GIT_FILEMODE_TREE)

    names = set(e.name for e in tree)
    names.add(subtree_name)
    try:
        names.remove("hash")
    except KeyError:
        pass

    m = hashlib.sha256()
    for n in sorted(names):
        if n == subtree_name:
            data = repo.get(subtree_oid)["hash"].read_raw()
        else:
            data = tree[n]["hash"].read_raw()
        m.update(data)

    treebuilder.insert("hash", repo.create_blob(m.digest()), git.GIT_FILEMODE_BLOB)
    return treebuilder.write()


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
            self._builder = self._repo.TreeBuilder(self._tree)

        else:
            self._commit = None
            self._repo = repo_create(repo_path)
            self._tree = self._repo.TreeBuilder().write()
            self._builder = self._repo.TreeBuilder(self._tree)
            self.write_commit(message="Init log")

    def add_leaf(self, input: str, content: bytes):
        sharded = shard(input)

        auto_insert(self._repo, self._builder, sharded, content)

        self._tree = self._builder.write()
        self.write_commit()

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
