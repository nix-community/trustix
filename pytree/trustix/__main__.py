from . import repo


if __name__ == '__main__':
    repo_path = "./repo"
    name = "trustix"
    email = "trustix@example.com"

    r = repo.Repository(repo_path, name=name, email=email)

    from random import choice
    from string import ascii_lowercase

    def rand_b(n) -> bytes:
        return ''.join(choice(ascii_lowercase) for i in range(n)).encode()

    for i in range(1):
        # Emulate hex-encoded input hash
        # TODO: Consider a different more compact encoding?
        input_hash = rand_b(64).decode()

        # Output NAR hash is assumed to be a binary sha256 hash
        # We're better off not serialising the value at all
        output_hash = rand_b(32)

        r.add_leaf(input_hash, output_hash)
