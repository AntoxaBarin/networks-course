import threading
from scapy.all import *
import time

total_in = 0
total_out = 0
monitor_ip = get_if_addr(conf.iface)
running = True


def packet_handler(packet):
    global total_in, total_out
    if IP in packet:
        size = len(packet)
        if packet[IP].src == monitor_ip:
            total_out += size
        else:
            total_in += size


def print_stats():
    while running:
        time.sleep(1)
        print(f"\rIncoming: {total_in} bytes | Outgoing: {total_out} bytes", end="")


threading.Thread(target=print_stats, daemon=True).start()

try:
    print(f"Monitoring traffic for IP: {monitor_ip}")
    sniff(prn=packet_handler, store=0)
except KeyboardInterrupt:
    running = False
    print("\n\nFinal stats:")
    print(f"Total incoming: {total_in} bytes")
    print(f"Total outgoing: {total_out} bytes")
