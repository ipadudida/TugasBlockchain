from base64 import encode
import sha3

kode = input("Masukan kode transaksi: ")
encoded = kode.encode()
obj_encoded = sha3.keccak_256(encoded)
print("kode transaksi setelah di hash:", obj_encoded.hexdigest())