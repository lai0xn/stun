package stunlib

import (
	"fmt"
	"net"
)


type ClientConfig struct {
  RemoteAddr string
}

type Client struct {
 cfg ClientConfig
}

func (client *Client) Dial(){
  header := Header {
    Type: BindingRequest,
    Length: 0,
    MagicCookie: magicCookie,
    TransactionID: [12]byte{1,2,3,4,5,6,7,8,9,10,11,12},
  }
  udpAddr,err := net.ResolveUDPAddr("udp4",client.cfg.RemoteAddr)
  if err != nil {
    panic(err)
  }
  encodedHeader := EncodeHeader(header)

  c,err := net.DialUDP("udp4",nil,udpAddr)
  if err != nil {
    panic(err)
  }

  _,err = c.Write(encodedHeader)
 
  if err != nil {
    panic(err)
  }
  buff := make([]byte,2048)
  n,_,err := c.ReadFromUDP(buff)
  if err != nil {
    panic(err)
  }
  fmt.Println(n)
  resHeader := DecodeHeader(buff)
  attrs := DecodeAttrs(buff[20:n],int(resHeader.Length))
  fmt.Println(resHeader.Length)
  fmt.Println(attrs)
}

func Dial(addr string){
  defaultClient := Client{
    cfg: ClientConfig{
      RemoteAddr: addr,
    },
  } 
  defaultClient.Dial()
}

