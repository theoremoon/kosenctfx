from math import gcd

exec(open("out.txt").read())

p = gcd(n, priv_key)
q = n // p

d = pow(e, -1, (p-1)*(q-1))
m = pow(c, d, n)

print(bytes.fromhex(hex(m)[2:]))

