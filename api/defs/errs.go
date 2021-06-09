package defs

//error handling
type Err struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}
type ErrResponse struct {
	HttpSC int
	Error  Err
}

var (
	//声明定义初始化了ErrResponse并复制给ErrRequestBodyParseFailed和ErrorAuthUser
	ErrorRequestBodyParseFailed = ErrResponse{HttpSC: 400, Error: Err{Error: "Request body is not correct", ErrorCode: "001"}}
	ErrorNotAuthUser            = ErrResponse{HttpSC: 401, Error: Err{Error: "User authentication failed", ErrorCode: "002"}}
	//创建用户时调用dbops插入数据库时错
	ErrorDBError = ErrResponse{HttpSC: 500, Error: Err{Error: "DB ops failed", ErrorCode: "003"}}
	//内部服务错误
	ErrorInternalFaults = ErrResponse{HttpSC: 500, Error: Err{Error: "Internal service error", ErrorCode: "004"}}
)
