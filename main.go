package main

import (
	"gopkg.in/yaml.v2"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"io"
	"net"
	"bufio"
	"os"
	"fmt"
	"log"
	"sync"
)

type ServerInfo struct {
	Ip        string `yaml:"ip"`
	Port      int    `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Path      string `yaml:"path"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usage: %s <config file> <cmd> \n", os.Args[0])
		return
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
		return
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
		return
	}
	servers := make([]ServerInfo, 0)
	err = yaml.Unmarshal(b, &servers)
	cmd := ""
	for i:=2; i<len(os.Args); i++ {
		cmd = cmd + " " + os.Args[i]
	}
	for _, s := range servers {
		addr := fmt.Sprintf("%v:%v", s.Ip, s.Port)
		prefix := s.Ip + ": "
		c := fmt.Sprintf("cd %v; %v", s.Path, cmd)
		stdout, stderr, _ := SSHCmd(s.Username, s.Password, addr, c)
		var wg sync.WaitGroup
		wg.Add(2)
		readFn := func(ss io.Reader) {
			defer wg.Done()
			reader := bufio.NewScanner(ss)
			for reader.Scan() {
				fmt.Print(prefix)
				fmt.Println(reader.Text())
			}
		}
		go readFn(stdout)
		go readFn(stderr)
		wg.Wait()
	}
}

func SSHCmd(user, password, addr, cmd string) (stdout,stderr io.Reader, err error) {
	log.Printf("%s %s %s %s", user, password, addr, cmd)
	passwd := []ssh.AuthMethod{ssh.Password(password)}
	conf := ssh.ClientConfig{
		User: user,
		Auth: passwd, 
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error{
		//log.Println("got hostKey:", hostname, key)
			return nil
		},
	}
	
	c, err := ssh.Dial("tcp", addr, &conf)
	if err != nil {
		return nil, nil, err
	}
	defer c.Close()
	if s, err := c.NewSession(); err != nil {
		return nil, nil, err
	} else {
		defer s.Close()
		stdout, _ = s.StdoutPipe()
		stderr, _ = s.StderrPipe()
		return stdout, stderr, s.Run(cmd)
	}
}

