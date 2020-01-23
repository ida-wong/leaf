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
	Server struct {
		MaxConnNum      int
		PendingWriteNum int
		MaxMsgLen       uint32
		Processor       network.Processor
		AgentChanRPC    *chanrpc.Server

		// websocket
		WSAddr       string
		HTTPTimeout  time.Duration
		CertFile     string
		KeyFile      string
		NewWsAgent   interface{}
		CloseWsAgent interface{}
		wsServer     *network.WSServer

		// tcp
		TCPAddr       string
		LenMsgLen     int
		LittleEndian  bool
		NewTcpAgent   interface{}
		CloseTcpAgent interface{}
		tcpServer     *network.TCPServer
	}

	serverAgent struct {
		conn     network.Conn
		server   *Server
		closeId  interface{}
		userData interface{}
	}
)

func (server *Server) Run(closeSig chan bool) {
	if server.WSAddr != "" {
		server.wsServer = new(network.WSServer)
		wsServer := server.wsServer
		wsServer.Addr = server.WSAddr
		wsServer.MaxConnNum = server.MaxConnNum
		wsServer.PendingWriteNum = server.PendingWriteNum
		wsServer.MaxMsgLen = server.MaxMsgLen
		wsServer.HTTPTimeout = server.HTTPTimeout
		wsServer.CertFile = server.CertFile
		wsServer.KeyFile = server.KeyFile
		if server.NewWsAgent == nil {
			server.NewWsAgent = "NewWsAgent"
		}
		if server.CloseWsAgent == nil {
			server.CloseWsAgent = "CloseWsAgent"
		}
		wsServer.NewAgent = func(conn *network.WSConn) network.Agent {
			agent := new(serverAgent)
			agent.conn = conn
			agent.server = server
			agent.closeId = server.CloseWsAgent
			if server.AgentChanRPC != nil {
				server.AgentChanRPC.Go(server.NewWsAgent, agent)
			}
			return agent
		}
	}

	if server.TCPAddr != "" {
		server.tcpServer = new(network.TCPServer)
		tcpServer := server.tcpServer
		tcpServer.Addr = server.TCPAddr
		tcpServer.MaxConnNum = server.MaxConnNum
		tcpServer.PendingWriteNum = server.PendingWriteNum
		tcpServer.LenMsgLen = server.LenMsgLen
		tcpServer.MaxMsgLen = server.MaxMsgLen
		tcpServer.LittleEndian = server.LittleEndian
		if server.NewTcpAgent == nil {
			server.NewTcpAgent = "NewTcpAgent"
		}
		if server.CloseTcpAgent == nil {
			server.CloseTcpAgent = "CloseTcpAgent"
		}
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			agent := new(serverAgent)
			agent.conn = conn
			agent.server = server
			agent.closeId = server.CloseTcpAgent
			if server.AgentChanRPC != nil {
				server.AgentChanRPC.Go(server.NewTcpAgent, agent)
			}
			return agent
		}
	}

	if server.wsServer != nil {
		server.wsServer.Start()
	}
	if server.tcpServer != nil {
		server.tcpServer.Start()
	}

	<-closeSig
}

func (server *Server) OnDestroy() {
	if server.wsServer != nil {
		server.wsServer.Close()
	}
	if server.tcpServer != nil {
		server.tcpServer.Close()
	}
}

func (agent *serverAgent) Run() {
	for {
		data, err := agent.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		if agent.server.Processor != nil {
			msg, err := agent.server.Processor.Unmarshal(data)
			if err != nil {
				log.Debug("unmarshal message error: %v", err)
				break
			}
			err = agent.server.Processor.Route(msg, agent)
			if err != nil {
				log.Debug("route message error: %v", err)
				break
			}
		}
	}
}

func (agent *serverAgent) OnClose() {
	if agent.server.AgentChanRPC != nil {
		err := agent.server.AgentChanRPC.Call0(agent.closeId, agent)
		if err != nil {
			log.Error("chanrpc error: %v", err)
		}
	}
}

func (agent *serverAgent) WriteMsg(msg interface{}) {
	if agent.server.Processor != nil {
		data, err := agent.server.Processor.Marshal(msg)
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

func (agent *serverAgent) LocalAddr() net.Addr {
	return agent.conn.LocalAddr()
}

func (agent *serverAgent) RemoteAddr() net.Addr {
	return agent.conn.RemoteAddr()
}

func (agent *serverAgent) Close() {
	agent.conn.Close()
}

func (agent *serverAgent) Destroy() {
	agent.conn.Destroy()
}

func (agent *serverAgent) UserData() interface{} {
	return agent.userData
}

func (agent *serverAgent) SetUserData(data interface{}) {
	agent.userData = data
}
