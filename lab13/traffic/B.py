import threading
from scapy.all import *
import time
from collections import defaultdict

total_in = 0
total_out = 0
monitor_ip = get_if_addr(conf.iface)
running = True

port_stats = defaultdict(lambda: {"in": 0, "out": 0})
last_report_time = time.time()


def packet_handler(packet):
    global total_in, total_out
    if IP in packet:
        size = len(packet)
        src_ip = packet[IP].src

        port = None
        if hasattr(packet, "sport") and hasattr(packet, "dport"):
            port = packet.dport if src_ip == monitor_ip else packet.sport

        if src_ip == monitor_ip:
            total_out += size
            if port is not None:
                port_stats[port]["out"] += size
        else:
            total_in += size
            if port is not None:
                port_stats[port]["in"] += size


def print_stats():
    global last_report_time
    while running:
        time.sleep(1)
        current_time = time.time()

        print(f"\rTotal - In: {total_in} bytes | Out: {total_out} bytes", end="")

        if current_time - last_report_time >= 5:
            last_report_time = current_time
            print("\n\n=== Traffic by port ===")

            sorted_ports = sorted(
                port_stats.items(), key=lambda x: x[1]["in"] + x[1]["out"], reverse=True
            )

            for port, stats in sorted_ports:
                total = stats["in"] + stats["out"]
                if total > 0:
                    print(
                        f"Port {port}: In={stats['in']} bytes | Out={stats['out']} bytes | Total={total} bytes"
                    )

            print("=" * 30 + "\n")


threading.Thread(target=print_stats, daemon=True).start()

try:
    print(f"Monitoring traffic for IP: {monitor_ip}")
    sniff(prn=packet_handler, store=0)
except KeyboardInterrupt:
    running = False
    print("\n\n=== Final Statistics ===")
    print(f"Total incoming: {total_in} bytes")
    print(f"Total outgoing: {total_out} bytes")

    if port_stats:
        print("\nPort traffic details:")
        for port, stats in sorted(
            port_stats.items(), key=lambda x: x[1]["in"] + x[1]["out"], reverse=True
        ):
            total = stats["in"] + stats["out"]
            if total > 0:
                print(
                    f"Port {port}: In={stats['in']} bytes | Out={stats['out']} bytes | Total={total} bytes"
                )
