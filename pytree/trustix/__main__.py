from . import repo
import os


if __name__ == '__main__':
    repo_path = "./repo"
    name = "trustix"
    email = "trustix@example.com"

    r = repo.Repository(repo_path, name=name, email=email)

    from random import choice
    from string import ascii_lowercase

    def rand_s(n):
        return ''.join(choice(ascii_lowercase) for i in range(n))

    for i in range(10 * 1000):
        # for i in range(1):
        print(i)
        input_hash = rand_s(32)  # Store hash
        output_hash = os.urandom(32)  # Output NAR hash
        r.add_leaf(input_hash, output_hash)

        new_input_hash = input_hash[:-4] + "aoeu"
        r.add_leaf(new_input_hash, output_hash)
