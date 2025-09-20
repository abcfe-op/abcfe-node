// cli 패키지는 명령행 인터페이스 (CLI)와 관련된 함수를 제공합니다.
package cli

import (
	"fmt"
	"os"

	"github.com/abcfe-op/abcfe-node/pos"
	"github.com/abcfe-op/abcfe-node/rest"
)

// cli 명령어 기본 가이드 (Ex. go run main.go -mode=rest -port=4000)
func usage() {
	fmt.Printf("Welcome to 민석's Blockchain Project\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-port:	Set the PORT of the server\n")
	fmt.Printf("-mode:	Choose between 'auto' and 'rest'\n")
	os.Exit(0)
}

// cli 명령어를 감지하여 auto 또는 rest 모드로 실행
func Start(port int, mode string) int {
	if len(os.Args) == 1 {
		usage()
	}

	switch mode {
	case "rest":
		rest.Start(port)
	case "auto":
		pos.PoS(port)
	default:
		usage()
	}

	return port
}
