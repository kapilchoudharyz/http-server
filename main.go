package main

import (
	"fmt"
	"strings"
	"syscall"
)

// This is a basic http server. This is to show you how http servers works under the hood.

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    interface{}
}

func ReqParser(req string) Request {
	fmt.Printf("Request received (raw): %q\n", req)
	if req == "" {
		fmt.Println("Empty request, returning empty slice")
		return Request{}
	}
	//time.Sleep(10 * time.Second)

	request := strings.Split(req, "\r\n")
	parsedRequest := Request{}
	headers := make(map[string]string)
	for i, line := range request {
		fmt.Printf("Line %d: %q\n", i, line)
		if i == 0 {
			S := strings.Split(line, " ")
			parsedRequest.Method = S[0]
			parsedRequest.Path = S[1]
			parsedRequest.Version = S[2]
		}
		if i > 0 && line != "" {
			S := strings.Split(line, ":")
			headers[S[0]] = S[1]
			parsedRequest.Headers = headers
		}
		if i > 0 && line == "" && i+1 < len(request) {
			parsedRequest.Body = request[i+1]
		}

	}
	fmt.Println("**********", parsedRequest)
	//if parsedRequest.Headers["Content-Length"] == "Application/json" {
	//
	//}

	return parsedRequest
}

func Server(port int) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		panic(err)
	}
	defer syscall.Close(fd)
	syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	addr := syscall.SockaddrInet4{Port: port, Addr: [4]byte{127, 0, 0, 1}}
	if err := syscall.Bind(fd, &addr); err != nil {
		panic(err)

	}
	syscall.Listen(fd, syscall.SOMAXCONN)
	for {
		connFd, _, err := syscall.Accept(fd)
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go func(cfd int) {
			defer syscall.Close(cfd)
			buf := make([]byte, 4096)
			n, err := syscall.Read(cfd, buf)
			if err != nil {
				fmt.Println("Read error:", err)
				return
			}

			req := string(buf[:n])
			request := ReqParser(req)
			fmt.Println("**********", request)
			//fmt.Print(request)
			//os.Stdout.Sync()

			response := "HTTP/1.1 200 OK\r\nContent-Length: 2\r\n\r\nOK"
			syscall.Write(cfd, []byte(response))
		}(connFd)
	}

}

func main() {
	Server(8080)
}
