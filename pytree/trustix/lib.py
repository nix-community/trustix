import pygit2 as git


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
