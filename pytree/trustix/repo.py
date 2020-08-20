import pygit2 as git  # type: ignore
import itertools
import string
import os.path
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
    depth = int(len(input) / 2)
    return tuple(
        input[2 * i : 2 * i + 2]
        for i in range(depth)
    )


def auto_insert(repo, treebuilder, path, thing, mode):
    print(path)
    path_parts = path.split('/', 1)
    if len(path_parts) == 1:
        treebuilder.insert(path, thing, mode)
        return treebuilder.write()

    subtree_name, sub_path = path_parts
    tree_oid = treebuilder.write()
    tree = repo.get(tree_oid)
    try:
        entry = tree[subtree_name]
        assert entry.filemode == git.GIT_FILEMODE_TREE,\
            '{} already exists as a blob, not a tree'.format(entry.name)
        existing_subtree = repo.get(entry.hex)
        sub_treebuilder = repo.TreeBuilder(existing_subtree)
    except KeyError:
        sub_treebuilder = repo.TreeBuilder()

    subtree_oid = auto_insert(repo, sub_treebuilder, sub_path, thing, mode)
    treebuilder.insert(subtree_name, subtree_oid, git.GIT_FILEMODE_TREE)
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

        else:
            self._commit = None
            self._repo = repo_create(repo_path)
            self._tree = self._repo.TreeBuilder().write()
            self.write_commit(message="Init log")
            # TODO: Create an empty root hash

        self._builder = self._repo.TreeBuilder(self._tree)

    def _write_tree(self, full_path: typing.Tuple[str, ...], content: bytes):
        """Write path to tree and update the internal tree reference"""

        oid = self._repo.write(git.GIT_OBJ_BLOB, content)

        visited = []
        for p in full_path[:-1]:
            visited.append(p)
            builder = self._repo.TreeBuilder(self._tree)
            builder.insert("/".join(visited), self._tree, git.GIT_FILEMODE_TREE)
            self._tree = builder.write()

        builder = self._repo.TreeBuilder(self._tree)
        builder.insert("/".join(full_path), oid, git.GIT_FILEMODE_BLOB)
        self._tree = builder.write()

    def add_leaf(self, input: str, content: bytes):
        sharded = os.path.sep.join(shard(input))

        auto_insert(self._repo, self._builder, sharded, self._repo.create_blob(content), git.GIT_FILEMODE_BLOB)
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
