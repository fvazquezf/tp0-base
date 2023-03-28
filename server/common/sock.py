# wrapper for send/recv of socket

def sendall(sock, data, size):
    total_sent = 0
    while total_sent < size:
        sent = sock.send(data[total_sent:])
        if sent == 0:
            raise RuntimeError("socket connection broken")
        total_sent += sent

def recvall(sock, size):
    data = bytearray(size)
    total_recv = 0
    while total_recv < size:
        recv_size = sock.recv_into(memoryview(data)[total_recv:])
        if recv_size == 0:
            raise RuntimeError("socket connection broken")
        total_recv += recv_size
    return data