package proxy

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type Proxy struct {
	Config string
}

func (p *Proxy) NewProxy(config string) *Proxy {
	return &Proxy{Config: config}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		p.HTTPSHandler(w, r)
	} else {
		p.HTTPHandler(w, r)
	}
}

func (p *Proxy) parseNewRequest(r *http.Request) Req {
	fmt.Println(r.Method)
	fmt.Println(r.Host)
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
	fmt.Println(r.Method)
	fmt.Println(r.URL)
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	p.CopyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
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

func (p *Proxy) HTTPSHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	fmt.Println(r.URL)

	destination, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		fmt.Println(ok)
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	client, _, err := hijacker.Hijack()
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	go p.TransferData(destination, client)
	go p.TransferData(client, destination)
}

func (p *Proxy) TransferData(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func (p *Proxy) CopyHeader(dest, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dest.Add(key, value)
		}
	}
}

type Req struct {
	method string
	host string
	path string
	userAgent string
	port string
}
