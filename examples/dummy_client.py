import sys

import requests


def send_request(address: str, length: int):
    inputs = {
        "length": length,
        "texts": ["Some text"] * length
    }
    resp = requests.post(address + "/infer", json=inputs)
    outputs = resp.json()
    total_time = outputs["process_time"]
    avg_time = total_time / length
    print(f"length {length}, total {total_time} time unit, average time {avg_time} / text")


if __name__ == "__main__":
    if len(sys.argv) < 2:
        raise ValueError
    address = sys.argv[1]
    lens = list(map(int, sys.argv[2:]))
    for l in lens:
        send_request(address, l)
