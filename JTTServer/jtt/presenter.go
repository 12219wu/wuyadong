package jtt

type Presenter interface {
	// 初始化
	Init(Context)
}

type BasePresenter struct {
	Ctx Context // 通信上下文
}

func (b *BasePresenter) Init(ctx Context) {
	b.Ctx = ctx
}

// JTT消息通信上下文
type Context interface {
	// 终端ID
	Client() Client

	// 终端消息
	Message() Input

	// 平台应答
	Response(Output)
}

// NewContext 新建协议通信上下文。
func NewContext(client Client, input Input) Context {
	return &contextImpl{
		client:  client,
		message: input,
	}
}

type contextImpl struct {
	client  Client // 终端
	message Input  // 终端消息
}

func (c *contextImpl) Client() Client {
	return c.client
}

func (c *contextImpl) Message() Input {
	return c.message
}

func (c *contextImpl) Response(output Output) {
	c.client.Send(output)
}
