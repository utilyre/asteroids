import struct

def send_message(sock, msg):
    send_integer(sock, msg["version"])
    send_string(sock, msg["entity"])
    send_string(sock, msg["action"])

def send_string(sock, value):
    encoded = value.encode("utf-8")
    send_integer(sock, len(encoded))
    sock.sendall(encoded)

def recv_string(sock):
    size = recv_integer(sock)
    data = sock.recv(size)
    return data.decode("utf-8")

def send_integer(sock, value):
    data = struct.pack(">Q", value)
    #                   ^ big-endian
    #                    ^ unsigned 8 bytes
    sock.sendall(data)

def recv_integer(sock):
    data = sock.recv(8)
    if len(data) < 8:
        raise ConnectionError("Incomplete data received")
    return struct.unpack(">Q", data)[0]
