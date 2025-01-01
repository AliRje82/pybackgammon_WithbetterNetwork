import socket
import base64
import rsa

private_key = []
public_key = []
for i in range(3):
    publicKey, privateKey = rsa.newkeys()
    private_key.append(privateKey)
    public_key.append(publicKey)

print("Keys are ready")

client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

client.connect(("localhost",8000))
message = b"Test"

for private,public in zip(private_key,public_key):
    client.send(private.save_pkcs1(format='PEM'))
    message=rsa.encrypt(message,public)

client.send(message)
