from . import repo


if __name__ == '__main__':
    repo_path = "./repo"
    name = "trustix"
    email = "trustix@example.com"

    r = repo.Repository(repo_path, name=name, email=email)

    from random import choice
    from string import ascii_lowercase

    def rand_s(n):
        return ''.join(choice(ascii_lowercase) for i in range(n))

    for i in range(1):
        print(i)
        r.add_leaf(rand_s(32), rand_s(64).encode())
