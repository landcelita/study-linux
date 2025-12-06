package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

const (
	ALLOC_SIZE = 1024 * 1024 * 1024
)

func main() {
	pid := os.Getpid()
	fmt.Println("*** 新規メモリ領域獲得前のメモリマップ ***")
	command := exec.Command("cat", "/proc/"+strconv.Itoa(pid)+"/maps")
	command.Stdout = os.Stdout
	if err := command.Run(); err != nil {
		log.Fatal(err)
	}

	data, err := syscall.Mmap(-1, 0, ALLOC_SIZE, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON|syscall.MAP_PRIVATE)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("")
	fmt.Println("*** 新規メモリ領域: アドレス = %p, サイズ = 0x%x ***\n", &data[0], ALLOC_SIZE)
	fmt.Println("")

	fmt.Println("*** 新規メモリ領域獲得後のメモリマップ ***")
	command = exec.Command("cat", "/proc/"+strconv.Itoa(pid)+"/maps")
	command.Stdout = os.Stdout
	if err := command.Run(); err != nil {
		log.Fatal(err)
	}
}
