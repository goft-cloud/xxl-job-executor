package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

var cmd *exec.Cmd

func main() {
	// 172.20.0.89:9999
	http.HandleFunc("/", index)
	http.HandleFunc("/beat", index)
	http.HandleFunc("/idleBeat", index)
	http.HandleFunc("/run", jobRun)
	http.HandleFunc("/kill", jobKill)
	http.HandleFunc("/log", jobLog)

	// 启动web服务，监听9090端口
	err := http.ListenAndServe("172.20.0.89:9996", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// w表示response对象，返回给客户端的内容都在对象里处理
// r表示客户端请求对象，包含了请求头，请求参数等等
func jobRun(w http.ResponseWriter, r *http.Request) {
	go func() {
		execCommand("/Users/stelin/goft/xjob/t.sh")
	}()

	// 往w里写入内容，就会在浏览器里输出
	fmt.Fprintf(w, "{\n    \"code\": 200,\n    \"msg\": null,\n    \"content\": null\n}")
}

// w表示response对象，返回给客户端的内容都在对象里处理
// r表示客户端请求对象，包含了请求头，请求参数等等
func jobKill(w http.ResponseWriter, r *http.Request) {
	pid := cmd.Process.Pid
	fmt.Println(pid)
	cmd.Process.Kill()
	// 往w里写入内容，就会在浏览器里输出
	fmt.Fprintf(w, "{\n    \"code\": 200,\n    \"msg\": null,\n    \"content\": null\n}")
}

// w表示response对象，返回给客户端的内容都在对象里处理
// r表示客户端请求对象，包含了请求头，请求参数等等
func jobLog(w http.ResponseWriter, r *http.Request) {
	s, _ := ioutil.ReadAll(r.Body) //把  body 内容读入字符串 s
	body := string(s)
	//{"logDateTim":1606479942000,"logId":75,"fromLineNum":1}
	fmt.Println(body)

	fmt.Println(r.URL.RequestURI())
	// 往w里写入内容，就会在浏览器里输出
	fmt.Fprintf(w, "{\n    \"code\": 200,\n    \"msg\": null,\n    \"content\": null\n}")
}

// w表示response对象，返回给客户端的内容都在对象里处理
// r表示客户端请求对象，包含了请求头，请求参数等等
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("stelin" + r.URL.Path)
	//time.Sleep(time.Duration(30) * time.Second)
	// 往w里写入内容，就会在浏览器里输出
	fmt.Fprintf(w, "{\n    \"code\": 200,\n    \"msg\": null,\n    \"content\": null\n}")
}

func execCommand(cd string) bool {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd = exec.Command("bash", "-c", cd)

	//显示运行的命令
	//StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return false
	}

	cmd.Start()
	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}

	//阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
	cmd.Wait()

	pid := cmd.Process.Pid
	fmt.Println("startPid=" + strconv.Itoa(pid))
	return true
}
