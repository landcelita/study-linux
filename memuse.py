#!/usr/bin/python3

import subprocess

size = 1_000_000

print("メモリ獲得前のシステム全体のメモリ使用量を表示します。")
subprocess.run("free")

array = [0] * size

print("メモリ獲得後のシステム全体のメモリ使用量を表示します。")
subprocess.run("free")