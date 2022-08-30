module netty_jtt808

go 1.15

require github.com/beego/beego/v2 v2.0.1

require (
	common/protocol v0.0.0-00010101000000-000000000000
	github.com/go-netty/go-netty v0.0.0-20220104093642-a83877336e91
	github.com/smartystreets/goconvey v1.6.4
)

replace common/protocol => ../common/protocol
