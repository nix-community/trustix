Because of limited support for partial repositories in [Go-git](https://github.com/go-git/go-git) this implementation is written in Python.
While this may have some impact for benchmarking the limiting factor should not be the language or Trustix itself but rather underlying cryptography & Git which is all implemented in native code.
Python is only used as a "glue language".
