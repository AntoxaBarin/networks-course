import socket
import struct
import time
import sys

TIMEOUT = 2
DEST = sys.argv[1]
PACKET_NUM = int(sys.argv[2]) if len(sys.argv) > 2 else 3
MAX_TTL = 30

sock = socket.socket(socket.AF_INET, socket.SOCK_RAW, socket.IPPROTO_ICMP)
sock.settimeout(2)
pack_num = 0


def checksum_calc(num) -> int:
    checksum = (num >> 16) + (num & 0xFFFF)
    checksum += checksum >> 16
    checksum = ~checksum & 0xFFFF
    return checksum


def create_pack() -> bytes:
    global pack_num
    pack_num += 1
    return struct.pack("!BBHHH", 8, 0, checksum_calc(2049 + pack_num), 1, pack_num)


def echo(packet, ttl, dest) -> bool:
    sock.setsockopt(socket.IPPROTO_IP, socket.IP_TTL, ttl)
    try:
        sock.sendto(packet, (socket.gethostbyname(dest), 0))
        start = time.time()
        data, address = sock.recvfrom(1024)
        rtt = "{:4.3f}".format((time.time() - start) * 1000)
        print(f"{ttl}.\t{rtt}\t", end=" ")
        try:
            name = socket.gethostbyaddr(address[0])
            print(f"{name[0]} ({address[0]})")
        except:
            print(f"{address[0]}")
        return data[20] == 0 & data[21] == 0
    except socket.timeout:
        print(f"    *", end="")
        return False


print(f"Destination {DEST}, TTL: {MAX_TTL}, Packets: {PACKET_NUM}")
dest_reached = False

for ttl in range(MAX_TTL):
    for _ in range(PACKET_NUM):
        dest_reached |= echo(create_pack(), ttl + 1, DEST)
    if dest_reached:
        break

sock.close()
