package conf

var (
	// 打印堆栈信息时的最大长度
	LenStackBuf = 4096

	// log
	LogLevel string // 日志等级
	LogPath  string // 日志文件所在目录
	LogFlag  int    // 日志前缀，参见标准库中的日志前缀配置

	// console
	ConsolePort   int               // 监听端口，用来处理自定义命令
	ConsolePrompt string = "Leaf# " //
	ProfilePath   string            // profile目录

	// cluster
	ListenAddr      string   //
	ConnAddrs       []string //
	PendingWriteNum int      //
)
