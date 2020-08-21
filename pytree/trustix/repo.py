import pygit2 as git  # type: ignore
import itertools
import string
import os.path
import hashlib
import typing
import time


# Precompute possible leaf permutations
_path_permutations: typing.List[str] = list("".join(p) for p in itertools.permutations(string.hexdigits.lower(), 2))


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
    # depth = int(len(input) / 2)
    depth = 5

    return tuple(
        input[2 * i : 2 * i + 2]
        for i in range(depth)
    ) + (input[2 * depth:],)


def insert_leaf(repo, treebuilder, path: typing.Tuple[str, ...], content: bytes):

    # We have reached the leaf, insert value
    if len(path) == 1:
        hash_contents = repo.create_blob(content)
        treebuilder.insert(path[0], hash_contents, git.GIT_FILEMODE_BLOB)
        return treebuilder.write()

    tree_oid = treebuilder.write()
    tree = repo.get(tree_oid)

    try:
        entry = tree[path[0]]
    except KeyError:
        treebuilder.insert(path[0], tree_oid, git.GIT_FILEMODE_TREE)
        tree = repo.get(treebuilder.write())
        entry = tree[path[0]]

    m = hashlib.sha256()
    for e in sorted(entry, key=lambda x: x.name):
        if e.name == "hash":
            continue

        content = b""

        if e.filemode == git.GIT_FILEMODE_BLOB:
            content = e.read_raw()
        else:  # Subtree
            try:
                h = e["hash"]
            except KeyError:
                continue

            content = h.read_raw()

        node = b":::".join((e.name.encode(), e.read_raw()))
        m.update(node)

    hash_contents = repo.create_blob(m.hexdigest())
    treebuilder.insert("hash", hash_contents, git.GIT_FILEMODE_BLOB)

    existing_subtree = repo.get(entry.hex)
    sub_treebuilder = repo.TreeBuilder(existing_subtree)

    subtree_oid = insert_leaf(repo, sub_treebuilder, path[1:], content)
    treebuilder.insert(path[0], subtree_oid, git.GIT_FILEMODE_TREE)
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

        insert_leaf(self._repo, self._builder, sharded, content)

        self._tree = self._builder.write()
        self.write_commit(input)

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
