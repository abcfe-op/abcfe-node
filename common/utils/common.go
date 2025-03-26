package utils

import (
	"os"
	"os/user"
	"path"
)

func HomeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

// 프로젝트 루트 디렉토리 찾기
func FindProjectRoot(startDir string) string {
	// 현재 디렉토리에서 시작해서 상위로 올라가며 go.mod 파일 찾기
	dir := startDir
	for {
		// go.mod 파일이 있는지 확인
		if _, err := os.Stat(path.Join(dir, "go.mod")); err == nil {
			return dir
		}

		// 상위 디렉토리로 이동
		parentDir := path.Dir(dir)
		if parentDir == dir {
			// 루트에 도달했으나 go.mod를 찾지 못함
			// 현재 디렉토리 반환
			return startDir
		}
		dir = parentDir
	}
}
