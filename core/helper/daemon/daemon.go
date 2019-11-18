package daemon

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/axetroy/go-fs"
)

type Action func() error

func getPidFilePath() (string, error) {
	executableNamePath, err := os.Executable()

	if err != nil {
		return "", err
	}

	return executableNamePath + ".pid", nil
}

func Start(action Action, shouldRunInDaemon bool) error {
	if shouldRunInDaemon && os.Getppid() != 1 {
		// 将命令行参数中执行文件路径转换成可用路径
		filePath, _ := filepath.Abs(os.Args[0])
		cmd := exec.Command(filePath, os.Args[1:]...)
		// 将其他命令传入生成出的进程
		// cmd.Stdin = os.Stdin // 给新进程设置文件描述符，可以重定向到文件中
		//cmd.Stdout = ioutil.Discard
		//cmd.Stderr = ioutil.Discard
		err := cmd.Start() // 开始执行新进程，不等待新进程退出
		return err
	} else {
		pidFilePath, err := getPidFilePath()

		if err != nil {
			return err
		}

		if err := fs.WriteFile(pidFilePath, []byte(fmt.Sprintf("%d", os.Getpid()))); err != nil {
			return err
		}

		return action()
	}
}

func Stop() error {
	pidFilePath, err := getPidFilePath()

	if err != nil {
		return err
	}

	if !fs.PathExists(pidFilePath) {
		return nil
	}

	b, err := fs.ReadFile(pidFilePath)

	if err != nil {
		return nil
	}

	pidStr := string(b)

	pid, err := strconv.Atoi(pidStr)

	if err != nil {
		return err
	}

	ps, err := os.FindProcess(pid)

	if err != nil {
		return err
	}

	if err := ps.Signal(syscall.SIGTERM); err != nil {
		return err
	}

	psState, err := ps.Wait()

	if err != nil {
		return err
	}

	haveBeenKill := psState.Exited()

	if haveBeenKill {
		log.Printf("进程 %d 已结束.\n", psState.Pid())

		_ = fs.Remove(pidFilePath)
	} else {
		log.Printf("进程 %d 结束失败.\n", psState.Pid())
	}

	return nil
}
