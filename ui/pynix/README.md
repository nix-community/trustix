## Pynix - Utility functions for working with Nix in Python

### Examples

#### Base32 encoding/decoding
``` python
import pynix

input = "v5sv61sszx301i0x6xysaqzla09nksnd"
b = pynix.b32decode(input)
output = pynix.b32encode(b)
assert input == output
```

#### Derivation parsing
``` python
import pynix

# Returns a dict with the same shape as nix show-derivation uses
d = pynix.drvparse(open("/nix/store/s6rn4jz1sin56rf4qj5b5v8jxjm32hlk-hello-2.10.drv").read())
assert isinstance(d, dict)
```
