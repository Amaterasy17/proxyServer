package proxy

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

var res = []byte("<html>\n<head><title>301 Moved Permanently</title></head>\n<body bgcolor=\"white\">\n<center><h1>301 Moved Permanently</h1></center>\n<hr><center>nginx/1.14.1</center>\n</body>\n</html>\n")

type Proxy struct {
	Config string
}



func (p *Proxy) NewProxy(config string) *Proxy {
	return &Proxy{Config: config}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newRequest := p.parseNewRequest(r)
	if newRequest.port == "" {
		newRequest.port = ":80"
	} else {
		newRequest.port = ""
	}
	r.RequestURI = ""
	r.Header.Del("Proxy-Connection")
	//newRequest := p.parseNewRequest(r)
	//result := p.MakeRequest(newRequest, r)
	//w.Write(result)
	conn, err := net.Dial("tcp", newRequest.host + newRequest.port)
	if err != nil {
		fmt.Println("error?")
	}

	reqData := newRequest.method + " " + newRequest.path + " HTTP/1.0\r\n" + "Host: " + newRequest.host + "\r\nUser-Agent: " + newRequest.userAgent + "\r\n\r\n"
	//reqData := r.Method + " " + r.URL.Path + " HTTP/1.1" + "\r\n\r\n"
	fmt.Println(reqData)
	//_, err = fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	 _, err = fmt.Fprintf(conn, reqData)
	if err != nil {
		fmt.Println("kak zaebal")
	}
	status := bufio.NewReader(conn)
	body, err := ioutil.ReadAll(status)
	if err != nil {
		fmt.Println("oi")
	}
	//fmt.Println(body)
	fmt.Printf("%s", body)
	w.Write(body)
	//p.HTTPHandler(w, r)
	//http.Client{}
}

func (p *Proxy) parseNewRequest(r *http.Request) Req {
	fmt.Println(r.Method)
	fmt.Println(r.Host)
	//fmt.Println(r.URL.Path)
	//fmt.Println(r.URL.Port())
	//fmt.Println(r.UserAgent())
	//fmt.Println(r.Body)
	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(result)
	return Req {
		method:    r.Method,
		host:      r.Host,
		path:      r.URL.Path,
		userAgent: r.UserAgent(),
		port: r.URL.Port(),
	}
}

func (p *Proxy) MakeRequest(req Req, r *http.Request) []byte {
	timeout := time.Duration(2 * time.Second)
	client := http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	//_, err := http.NewRequest("GET", "http://" + req.host, nil)
	//if err != nil {
	//	fmt.Println(err)
	//}
	r.RequestURI = ""
	r.Header.Del("Proxy-Connection")
	resp, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}


	result2, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result2)
	return result2
}

func (p *Proxy) HTTPHandler(w http.ResponseWriter, r *http.Request)  {
	newRequest := p.parseNewRequest(r)
	if newRequest.port == "" {
		newRequest.port = ":80"
	} else {
		newRequest.port = ""
	}
	r.RequestURI = ""
	r.Header.Del("Proxy-Connection")

	conn, err := net.Dial("tcp", newRequest.host + newRequest.port)
	if err != nil {
		log.Println(err)
		return
	}

	reqData := newRequest.method + " " + newRequest.path + " HTTP/1.1\n" + "Host: " + newRequest.host + "\nUser-Agent: " + newRequest.userAgent + "\n\n"
	fmt.Println(reqData)
	fmt.Fprintf(conn, reqData)

	//_, err = conn.Write([]byte(reqData))
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	w.Write([]byte("error"))
	//	return
	//}
	defer conn.Close()

	status := bufio.NewReader(conn)
	body, err := ioutil.ReadAll(status)
	if err != nil {
		fmt.Println("oi")
	}
	fmt.Printf("%s", body)

	//var readChan chan []byte = make(chan []byte)
	//p.ReadResponse(conn, readChan)
	//fmt.Println("After read response")
	//var reader io.Reader
	//var read []byte
	//var readChan chan []byte = make(chan []byte)
	//n, err := conn.Read(read)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	w.Write([]byte("error"))
	//	return
	//}
	//if n != 0 {
	//	readChan <- read
	//} else {
	//	fmt.Println("fuck")
	//	fmt.Println(<-readChan)
	//}

	//conn.SetDeadline(time.Now().Add(time.Millisecond * 100))
	w.WriteHeader(http.StatusOK)
	//io.Copy(w, strings.NewReader(string(<-readChan)))
	//io.Copy(w, conn)
	w.Write(body)

}

func (p *Proxy) ReadResponse(conn net.Conn, readChan chan []byte) {
	fmt.Println("I in read response")
	var read []byte
	for {
		n, err := conn.Read(read)
		if err != nil {
			close(readChan)
			break
		}
		if n != 0 {
			fmt.Println("Reading")
		}
		readChan <- read
	}
	fmt.Println("I in end response")
}

type Req struct {
	method string
	host string
	path string
	userAgent string
	port string
}
