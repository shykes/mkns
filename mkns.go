package main

import (
	"fmt"
	"os"
	"path"
	"syscall"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s net|pid|mnt PATH\n", os.Args[0])
		os.Exit(1)
	}
	if err := mkns(os.Args[1], os.Args[2]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Println(os.Args[2])
}

func mkns(nsType, mntPath string) error {
	flags := map[string]int{
		"net": syscall.CLONE_NEWNET,
		"mnt": syscall.CLONE_NEWNS,
		"pid": syscall.CLONE_NEWPID,
		"ipc": syscall.CLONE_NEWIPC,
	}
	flag, ok := flags[nsType]
	if !ok {
		return fmt.Errorf("unsupported namespace type: %s", nsType)
	}
	if err := syscall.Unshare(flag); err != nil {
		return err
	}
	dir, _ := path.Split(mntPath)
	f, err := os.OpenFile(mntPath, os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		return err
	}
	f.Close()
	_ = os.MkdirAll(dir, 0700)
	if err := syscall.Mount("/proc/self/ns/net", mntPath, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return err
	}
	return nil
}
