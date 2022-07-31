//go:build ignore

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/rajveermalviya/gamen/scripts/utils"
)

var goroot = func() string {
	out, err := exec.Command("go", "env", "GOROOT").Output()
	must(err)

	var sb strings.Builder
	sb.Write(out)
	return strings.TrimSpace(sb.String())
}()

func main() {
	fmt.Println("rm -rf _web/")
	must(os.RemoveAll("_web/"))

	cmd := exec.Command("go", "build", "-o", "_web/test.wasm", "github.com/rajveermalviya/gamen/examples/hello")
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	fmt.Println(cmd.String())
	must(cmd.Run())

	fmt.Println("cp", goroot+"/misc/wasm/wasm_exec.js", "_web/wasm_exec.js")
	must(utils.Cp(goroot+"/misc/wasm/wasm_exec.js", "_web/wasm_exec.js"))

	fmt.Println("cp", goroot+"/misc/wasm/wasm_exec.html", "_web/index.html")
	must(utils.Cp(goroot+"/misc/wasm/wasm_exec.html", "_web/index.html"))

	fmt.Println("cd _web/")
	must(os.Chdir("_web/"))

	fmt.Println("serving _web/ directory on :8080")
	panic(http.ListenAndServe(":8080", http.FileServer(http.FS(os.DirFS(".")))))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
