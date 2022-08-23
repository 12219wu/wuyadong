package jtt

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/transport"
	"github.com/go-netty/go-netty/transport/tcp"
)

var (
	JttApp Server
)

func init() {
	JttApp = newJttServer()
}

// Run 运行jtt服务
//
// addr 服务运行地址
//
// jtt.Run("127.0.0.1:8081")
func Run(addr string) {
	tcpOptions := &tcp.Options{
		Timeout:         time.Second * 3,
		KeepAlive:       true,
		KeepAlivePeriod: time.Second * 5,
		Linger:          0,
		NoDelay:         true,
		SockBuf:         2048,
	}

	JttApp.listenAsync(addr, tcp.WithOptions(tcpOptions))
}

// newJttServer 创建JTT部标服务器
func newJttServer() Server {
	pipelineInitializer := func(channel netty.Channel) {
		channel.Pipeline().
			AddLast(DelimiterCodec(0x7E, true, 2048)).
			AddLast(EscapeCodec(0x7D, map[byte]byte{0x7D: 0x01, 0x7E: 0x02}, 0x7D, map[byte]byte{0x01: 0x7D, 0x02: 0x7E})).
			AddLast(PacketCodec()).AddLast(MessageCodec(0, nil, sequenceCounter()))
	}

	server := &server{
		channelIDFactory:  netty.SequenceID(),
		pipelineFactory:   netty.NewPipeline(),
		channelFactory:    netty.NewChannel(128),
		transportFactory:  tcp.New(),
		clientInitializer: pipelineInitializer,
		routes:            make(map[uint16]*route),
	}
	server.ctx, server.cancel = context.WithCancel(context.Background())

	return server
}

// Router 添加一条协议路线到JttApp中。
// 用法：
//  jtt.Router(MsgIdTerminalLogin, &LoginPresenter{}, "TerminalLogin")
//  jtt.Router(MsgIdTerminalAuth, &LoginPresenter{}, "TerminalAuth")
func Router(msgId uint16, p Presenter, methodName string) {
	JttApp.router(msgId, p, methodName)
}

// route 协议通信路线。
// 通过向路由器中注册终端主发消息的处理函数，在
// 接收到指定消息时，路由器会调用指定处理函数。
type route struct {
	msgId         uint16           // 消息ID
	presenterType reflect.Type     // 控制器
	methodName    string           // 执行函数名
	initialize    func() Presenter // 控制器初始化函数
}

// Server JTT部标服务接口
type Server interface {
	// 获取终端
	GetClient(id uint32) Client

	listen(url string, option ...transport.Option) error

	listenAsync(url string, option ...transport.Option)

	router(msgId uint16, p Presenter, methodName string)
}

// server JTT部标服务
type server struct {
	ctx               context.Context
	cancel            context.CancelFunc
	clientInitializer netty.ChannelInitializer
	transportFactory  netty.TransportFactory
	channelFactory    netty.ChannelFactory
	pipelineFactory   netty.PipelineFactory
	channelIDFactory  netty.ChannelIDFactory
	acceptor          transport.Acceptor
	clients           sync.Map
	routes            map[uint16]*route
}

// serveTransport to serve channel
func (bs *server) serveTransport(transport transport.Transport) Client {

	// create a new pipeline
	pipeline := bs.pipelineFactory()

	// generate a channel id
	cid := bs.channelIDFactory()

	// create a channel
	channel := bs.channelFactory(cid, bs.ctx, pipeline, transport)

	// create a client
	client := &client{
		id:        generateHash(transport.RemoteAddr().String()),
		channel:   channel,
		subsriber: bs,
	}

	// set the attachment if necessary
	channel.SetAttachment(client)

	// initialization pipeline
	bs.clientInitializer(channel)

	// serve channel.
	channel.Pipeline().ServeChannel(channel)
	return client
}

func (s *server) listen(url string, option ...transport.Option) error {
	if nil != s.acceptor {
		return fmt.Errorf("duplicate call Listener:Sync")
	}

	var err error
	var options *transport.Options
	if options, err = transport.ParseOptions(s.ctx, url, option...); nil != err {
		return err
	}

	if s.acceptor, err = s.transportFactory.Listen(options); nil != err {
		return err
	}

	for {
		// accept the transport
		t, err := s.acceptor.Accept()
		if nil != err {
			return err
		}

		select {
		case <-s.ctx.Done():
			// bootstrap has been closed
			return t.Close()
		default:
			// serve child transport
			log.Printf("终端[%s]已接入", t.RemoteAddr())
			s.serveTransport(t)
		}
	}
}

func (s *server) listenAsync(url string, option ...transport.Option) {
	go func() {
		if err := s.listen(url, option...); nil != err &&
			!strings.Contains(err.Error(), "use of closed network connection") {
			log.Fatal(err)
		}
	}()
}

func (s *server) GetClient(id uint32) Client {
	if value, ok := s.clients.Load(id); ok {
		return value.(Client)
	}
	return nil
}

func (s *server) onClientConnected(client Client) {
	s.clients.Store(client.ID(), client)
}

func (s *server) onClientDisconnected(client Client) {
	s.clients.Delete(client.ID())
}

func (s *server) onMessage(client Client, input Input) {
	route := s.routes[input.msgID()]
	if nil == route {
		return
	}

	ctx := NewContext(client, input)

	presenter := route.initialize()
	presenter.Init(ctx)
	method := reflect.ValueOf(presenter).MethodByName(route.methodName)

	go func() {
		method.Call(nil)
	}()
}

func (s *server) router(msgId uint16, p Presenter, methodName string) {
	route := route{
		msgId:      msgId,
		methodName: methodName,
	}

	impl := reflect.ValueOf(p)
	t := reflect.Indirect(impl).Type()
	method := impl.MethodByName(methodName)
	if !method.IsValid() {
		log.Panicf("方法【%s】不在类型【%s】中。", methodName, t.Name())
	}

	route.presenterType = t
	route.initialize = func() Presenter {
		v := reflect.New(route.presenterType)
		execPresenter, ok := v.Interface().(Presenter)
		if !ok {
			log.Println("非Presenter接口")
		}

		elemVal := reflect.ValueOf(p).Elem()
		elemType := reflect.TypeOf(p).Elem()
		execElem := reflect.ValueOf(execPresenter).Elem()

		numOfFields := elemVal.NumField()
		for i := 0; i < numOfFields; i++ {
			fieldType := elemType.Field(i)
			elemField := execElem.FieldByName(fieldType.Name)
			if elemField.CanSet() {
				fieldVal := elemVal.Field(i)
				elemField.Set(fieldVal)
			}
		}

		return execPresenter
	}

	s.routes[msgId] = &route
}
