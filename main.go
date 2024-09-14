package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
	"time"
)

func main() {
	// monitor()

	// 创建或打开一个文件用于写入崩溃报告
	crashFile, err := os.OpenFile("crash_report.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("无法打开崩溃报告文件:", err)
		return
	}
	defer crashFile.Close()

	debug.SetCrashOutput(crashFile, debug.CrashOptions{})

	// 设置一个 defer 函数来捕获整个程序中的 panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("程序发生严重错误: %v", r)
			// 这里可以添加一些清理代码
			// 例如关闭文件、数据库连接等
		}
	}()

	fmt.Println("程序开始运行")

	// 模拟一些可能导致 panic 的操作
	go dangerousOperation()
	time.Sleep(1 * time.Second)

	fmt.Println("程序正常结束") // 如果发生 panic，这行不会被执行
}

func dangerousOperation() {
	// 模拟一个导致 panic 的操作
	panic("发生了一个严重错误")
}
func monitor() {
	const monitorVar = "RUNTIME_DEBUG_MONITOR"
	if os.Getenv(monitorVar) != "" {
		// 实际演示 debug.SetCrashOutput 设置后的逻辑
		log.SetFlags(0)
		log.SetPrefix("monitor: ")

		crash, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("failed to read from input pipe: %v", err)
		}
		if len(crash) == 0 {
			os.Exit(0)
		}

		f, err := os.CreateTemp("", "*.crash")
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write(crash); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
		log.Fatalf("saved crash report at %s", f.Name())
	}

	// 模拟应用程序进程，设置 debug.SetCrashOutput 值
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(exe, "-test.run=ExampleSetCrashOutput_monitor")
	cmd.Env = append(os.Environ(), monitorVar+"=1")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stderr
	pipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("StdinPipe: %v", err)
	}
	debug.SetCrashOutput(pipe.(*os.File), debug.CrashOptions{})
	if err := cmd.Start(); err != nil {
		log.Fatalf("can't start monitor: %v", err)
	}

}
