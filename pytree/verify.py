import pygit2 as git
import subprocess
import hashlib
from trustix.repo import shard, format_hash


# git clone --filter=blob:none --no-checkout --sparse file://(pwd)/repo r
# git clone --filter=blob:none --bare file://(pwd)/repo r

repo_path = "r"


repo = git.Repository(repo_path)
commit = repo.head.resolve().target
tree = repo.get(commit).tree.id


def fetch_oid(oid: str):
    # Until https://github.com/libgit2/libgit2/pull/5603 is resolved we do this by shelling out to git
    # As it is right now libgit2 doesn't allow fetching specific OIDs
    #
    # It's possible that go-git is better at this and that we can use that instead once we've ported from python
    subprocess.run(
        ["git", "fetch", "origin", oid],
        cwd=repo_path,
        check=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )


def find_tree_obj(tree, path):
    x = tree[path[0]]
    if x.name == path[0]:
        if x.filemode == git.GIT_FILEMODE_BLOB:
            return x
        elif x.filemode == git.GIT_FILEMODE_TREE:
            return find_tree_obj(x, path[1:])
        else:
            raise ValueError(f"Cannot handle filemode '{x.filemode}'")

    raise ValueError(f"Could not find {path[0]} in tree {tree}")


def get_oid(oid):
    val = repo.get(oid)
    if val is None:
        fetch_oid(str(oid))
    return repo.get(oid)


def hash_tree(root_tree, path):
    """Hash a tree (used for verification)"""

    tree_sub = root_tree[path[0]]

    m = hashlib.sha256()

    for sub in sorted(tree_sub, key=lambda x: x.name):
        if sub.name == "hash":
            continue

        if sub.filemode == git.GIT_FILEMODE_TREE and sub.name == path[1]:
            data = hash_tree(tree_sub, path[1:])
        elif sub.filemode == git.GIT_FILEMODE_TREE:
            data = get_oid(sub["hash"].oid).read_raw()
        else:
            data = get_oid(sub.oid).read_raw()

        m.update(format_hash(sub.name, data))

    return m.digest()


def get_leaf(leaf: str):
    root_tree = repo.get(tree)
    path = shard(leaf)

    aggregate_hash = hash_tree(root_tree, path[:-1])

    m = hashlib.sha256()

    # Now read the rest of the hashes
    for sub in sorted(root_tree, key=lambda x: x.name):
        if sub.filemode != git.GIT_FILEMODE_TREE:
            continue

        if sub.name == path[0]:
            value = aggregate_hash
        else:
            value = get_oid(sub["hash"].oid).read_raw()

        m.update(format_hash(sub.name, value))

    if m.digest() != get_oid(root_tree["hash"].oid).read_raw():
        raise ValueError("Root hash mismatch")

    return


get_leaf("lnskpntaivodjqlvesbynelyholaaoeu")
get_leaf("phgvuusitqrxgxpyinvozomsqlgdaoeu")
