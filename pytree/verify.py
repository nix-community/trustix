import pygit2 as git
import subprocess
from trustix.repo import shard


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
    for x in tree:
        if x.name == path[0]:
            if x.filemode == git.GIT_FILEMODE_BLOB:
                return x
            return find_tree_obj(x, path[1:])
    raise ValueError(f"Could not find {path[0]} in tree {tree}")


def verify_leaf(leaf: str):
    path = shard(leaf)
    blob = find_tree_obj(repo.get(tree), path)
    fetch_oid(str(blob.oid))

    output_hash: bytes = blob.read_raw()

    print(output_hash)


verify_leaf("wasdszkjhdznunvavamlliviufdxsfeg")
