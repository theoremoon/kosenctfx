from flag import flag
from Crypto.Util.number import getStrongPrime

e = 65537

p = getStrongPrime(1024)
q = getStrongPrime(1024)
n = p * q


c = pow(int.from_bytes(flag, "big"), e, n)
priv_key = pow(p, e, n)

print(f"{n=}")
print(f"{e=}")
print(f"{c=}")
print(f"{priv_key=}")
