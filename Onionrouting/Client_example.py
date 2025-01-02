from Crypto.Cipher import AES
from Crypto.Random import get_random_bytes
from Crypto.Util.Padding import pad, unpad
import socket
import struct

#Encrypt
def encrypt_aes_ecb(key, plaintext):
    cipher = AES.new(key, AES.MODE_ECB)  
    padded_data = pad(plaintext, AES.block_size)  
    ciphertext = cipher.encrypt(padded_data)
    return ciphertext

# Decrypt
def decrypt_aes_ecb(key, ciphertext):
    try:
        cipher = AES.new(key, AES.MODE_ECB)
        decrypted_data = cipher.decrypt(ciphertext)
        print(f"Decrypted (raw): {decrypted_data}")
        return decrypted_data
    except Exception as e:
        print(f"Decryption error: {e}")
        raise

#Form a packet with header!
def make_pkt(message):
    message_length = len(message)
    length_header = struct.pack("!I", message_length)
    return length_header + message

def receive_message(sock):
    length_header = sock.recv(4)
    if not length_header:
        return None
    message_length = struct.unpack("!I", length_header)[0]

    message = b""
    while len(message) < message_length:
        chunk = sock.recv(message_length - len(message))
        if not chunk:
            raise ValueError("Socket connection broken")
        message += chunk
    return message


def key_gen():
    keys = []
    messages = []

    for i in range(3):
        key = get_random_bytes(32)
        message = key
        for j in range(len(keys)-1, -1, -1):
            message = encrypt_aes_ecb(keys[j], message)
        messages.append(make_pkt(message))
        keys.append(key)
    return keys,messages

#Encrypt message
def encrypt_message(keys,message):
    for k in reversed(keys):
        message = encrypt_aes_ecb(k,message)
    return message

#Decrypt the message
def decrypt_message(keys, message):
    for k in keys:
        message = decrypt_aes_ecb(k, message)
        try:
            
            message = unpad(message, AES.block_size)
        except ValueError:
            pass
    return message


#Transform the keys!
def exchange_keys(client,messages):
    for m in messages:
        client.send(m)


#Dice rolling
def dice_rolls(client,message,keys):
    client.send(make_pkt(encrypt_message(keys,message)))
    message = receive_message(client)
    message = decrypt_message(keys,message)
    message = message.split(b',')
    try:
        dice_1 = int(message[0][0])
        dice_2 = int(message[1][0])
        return (dice_1,dice_2)
    except Exception as e:
        return -1



keys,exchange_messages=key_gen()
client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
client.connect(("localhost", 8000))
exchange_keys(client,exchange_messages)
loging_message = b"ali,localhost:9000"

client.send(make_pkt(encrypt_message(keys,loging_message)))

client.close()



