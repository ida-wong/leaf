package gate

import (
	"net"
	"reflect"
	"time"

	"github.com/ida-wong/leaf/chanrpc"
	"github.com/ida-wong/leaf/log"
	"github.com/ida-wong/leaf/network"
)

type (
	Client struct {
		PendingWriteNum int
		MaxMsgLen       uint32
		Processor       network.Processor
		AgentChanRPC    *chanrpc.Server
		AutoReconnect   bool
		ConnectInterval time.Duration

		// websocket
		WSAddr           string
		WsConnNum        int
		HandshakeTimeout time.Duration
		NewWsAgent       interface{}
		CloseWsAgent     interface{}
		wsClient         *network.WSClient

		// tcp
		TCPAddr       string
		LenMsgLen     int
		LittleEndian  bool
		TcpConnNum    int
		NewTcpAgent   interface{}
		CloseTcpAgent interface{}
		tcpClient     *network.TCPClient
	}

	clientAgent struct {
		conn     network.Conn
		client   *Client
		closeId  interface{}
		userData interface{}
	}
)

func (client *Client) Run(closeSig chan bool) {
	if client.WSAddr != "" {
		client.wsClient = new(network.WSClient)
		wsClient := client.wsClient
		wsClient.Addr = client.WSAddr
		wsClient.ConnNum = client.WsConnNum
		wsClient.ConnectInterval = client.ConnectInterval
		wsClient.PendingWriteNum = client.PendingWriteNum
		wsClient.MaxMsgLen = client.MaxMsgLen
		wsClient.HandshakeTimeout = client.HandshakeTimeout
		wsClient.AutoReconnect = client.AutoReconnect
		if client.NewWsAgent == nil {
			client.NewWsAgent = "NewWsAgent"
		}
		if client.CloseWsAgent == nil {
			client.CloseWsAgent = "CloseWsAgent"
		}
		wsClient.NewAgent = func(conn *network.WSConn) network.Agent {
			agent := new(clientAgent)
			agent.conn = conn
			agent.client = client
			agent.closeId = client.CloseWsAgent
			if client.AgentChanRPC != nil {
				client.AgentChanRPC.Go(client.NewWsAgent, agent)
			}
			return agent
		}

		wsClient.Start()
	}

	if client.TCPAddr != "" {
		client.tcpClient = new(network.TCPClient)
		tcpClient := client.tcpClient
		tcpClient.Addr = client.TCPAddr
		tcpClient.ConnNum = client.TcpConnNum
		tcpClient.ConnectInterval = client.ConnectInterval
		tcpClient.PendingWriteNum = client.PendingWriteNum
		tcpClient.AutoReconnect = client.AutoReconnect
		tcpClient.LenMsgLen = client.LenMsgLen
		tcpClient.MaxMsgLen = client.MaxMsgLen
		tcpClient.LittleEndian = client.LittleEndian
		if client.NewTcpAgent == nil {
			client.NewTcpAgent = "NewTcpAgent"
		}
		if client.CloseTcpAgent == nil {
			client.CloseTcpAgent = "CloseTcpAgent"
		}
		tcpClient.NewAgent = func(conn *network.TCPConn) network.Agent {
			agent := new(clientAgent)
			agent.conn = conn
			agent.client = client
			agent.closeId = client.CloseTcpAgent
			if client.AgentChanRPC != nil {
				client.AgentChanRPC.Go(client.NewTcpAgent, agent)
			}
			return agent
		}

		tcpClient.Start()
	}

	<-closeSig
}

func (client *Client) OnDestroy() {
	if client.wsClient != nil {
		client.wsClient.Close()
	}
	if client.tcpClient != nil {
		client.tcpClient.Close()
	}
}

func (agent *clientAgent) Run() {
	for {
		data, err := agent.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		if agent.client.Processor != nil {
			msg, err := agent.client.Processor.Unmarshal(data)
			if err != nil {
				log.Debug("unmarshal message error: %v", err)
				break
			}
			err = agent.client.Processor.Route(msg, agent)
			if err != nil {
				log.Debug("route message error: %v", err)
				break
			}
		}
	}
}

func (agent *clientAgent) OnClose() {
	if agent.client.AgentChanRPC != nil {
		err := agent.client.AgentChanRPC.Call0(agent.closeId, agent)
		if err != nil {
			log.Error("chanrpc error: %v", err)
		}
	}
}

func (agent *clientAgent) WriteMsg(msg interface{}) {
	if agent.client.Processor != nil {
		data, err := agent.client.Processor.Marshal(msg)
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		err = agent.conn.WriteMsg(data...)
		if err != nil {
			log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
		}
	}
}

func (agent *clientAgent) LocalAddr() net.Addr {
	return agent.conn.LocalAddr()
}

func (agent *clientAgent) RemoteAddr() net.Addr {
	return agent.conn.RemoteAddr()
}

func (agent *clientAgent) Close() {
	agent.conn.Close()
}

func (agent *clientAgent) Destroy() {
	agent.conn.Destroy()
}

func (agent *clientAgent) UserData() interface{} {
	return agent.userData
}

func (agent *clientAgent) SetUserData(data interface{}) {
	agent.userData = data
}
