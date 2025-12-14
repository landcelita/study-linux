/*
cache

1. 2**2(4)KBから、2**4.25KB, 2**4.5KB, ...と最終的には64MBの数値に対して以下の処理をする
 1. 数値に相当するサイズのバッファを獲得
 2. バッファの全キャッシュラインにシーケンシャルにアクセス。最後のキャッシュラインへのアクセスが終わったら最初のキャッシュラインに戻り、最終的にはソースコードに書いてあるNACCESS回メモリアクセスする
 3. 1回アクセスあたりの所要時間を記録

2. 1の結果を元にcache.jpgというファイルにグラフを出力
*/
package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"syscall"
	"time"
)

const (
	CacheLineSize = 64
	NACCESS       = 128 * 1024 * 1024
)

func main() {
	_ = os.Remove("out.txt")
	f, err := os.OpenFile("out.txt", os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		log.Fatalf("Failed to open out.txt file: %v", err)
	}
	defer f.Close()

	for i := 2.0; i <= 22.5; i += 0.25 {
		bufSize := int(math.Pow(2, i)) * 1024
		data, err := syscall.Mmap(-1, 0, bufSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON|syscall.MAP_PRIVATE)
		defer syscall.Munmap(data)
		if err != nil {
			log.Fatalf("Failed to mmap: %v", err)
		}

		fmt.Printf("バッファサイズ 2**%.2f (%d) KB についてのデータを収集中...\n", i, bufSize/1024)
		start := time.Now()
		for i := 0; i < NACCESS/(bufSize/CacheLineSize); i++ {
			for j := 0; j < bufSize; j += CacheLineSize {
				data[j] = 0
			}
		}
		end := time.Since(start)
		f.Write([]byte(fmt.Sprintf("%f\t%f\n", i, float64(NACCESS)/float64(end.Nanoseconds()))))
	}
	command := exec.Command("python3 plot-cache.py")
	out, err := command.Output()
	if err != nil {
		log.Fatalf("Failed to execute plot-cache.py: %v. out: %v", err, string(out))
		os.Exit(1)
	}
}
