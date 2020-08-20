from . import repo


if __name__ == '__main__':
    repo_path = "./repo"
    name = "trustix"
    email = "trustix@example.com"

    r = repo.Repository(repo_path, name=name, email=email)

    r.add_leaf("w9yy7v61ipb5rx6i35zq1mvc2iqfmps1", b"sha256:1mi14cqk363wv368ffiiy01knardmnlyphi6h9xv6dkjz44hk30i")
    r.add_leaf("w9yy7v61ipb5rx6i35zq1mvc2iqfmpbb", b"sha256:1mi14cqk363wv368ffiiy01knardmnlyphi6h9xv6dkjz44hk30a")
