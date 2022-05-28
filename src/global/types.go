package types

// 说明：
// 1. 所提到的「位数」均以字节长度为准
// 2. 所有的 ID 均为 int64（以 string 方式表现）

// 通用结构

type ErrNo int

const (
	OK             ErrNo = 0
	ParamInvalid   ErrNo = 1   // 参数不合法
	UserHasExisted ErrNo = 2   // 该 Username 已存在
	UserHasDeleted ErrNo = 3   // 用户已删除
	UserNotExisted ErrNo = 4   // 用户不存在
	WrongPassword  ErrNo = 5   // 密码错误
	LoginRequired  ErrNo = 6   // 用户未登录
	UnknownError   ErrNo = 255 // 未知错误
)

type ResponseMeta struct {
	Code ErrNo
}

type TMember struct {
	UserID   string
	Nickname string
	Username string
	hobby    string
	UserType UserType
}

type TData struct {
	DataID    string
	Content   string
	ImagesUrl string
}

// 成员管理

type UserType int

const (
	Admin  UserType = 1
	Normal UserType = 2
)

// 系统内置管理员账号
// 账号名：JudgeAdmin 密码：JudgePassword2022

// 1.创建成员
// 参数不合法返回 ParamInvalid

// 只有管理员才能添加

type CreateMemberRequest struct {
	Nickname string
	Username string
	Password string
	UserType UserType
}

type CreateMemberResponse struct {
	Code ErrNo
	Data struct {
		UserID string // int64 范围
	}
}

// 2.获取成员信息

type GetMemberRequest struct {
	UserID string
}

// 如果用户已删除请返回已删除状态码，不存在请返回不存在状态码

type GetMemberResponse struct {
	Code ErrNo
	Data TMember
}

// 3.批量获取成员信息

type GetMemberListRequest struct {
	Offset int
	Limit  int
}

type GetMemberListResponse struct {
	Code ErrNo
	Data struct {
		MemberList []TMember
	}
}

// 4.更新成员信息

type UpdateMemberRequest struct {
	UserID   string
	Nickname string
}

type UpdateMemberResponse struct {
	Code ErrNo
}

// 5.删除成员信息
// 成员删除后，该成员不能够被登录且不应该不可见，ID 不可复用

type DeleteMemberRequest struct {
	UserID string
}

type DeleteMemberResponse struct {
	Code ErrNo
}

// ----------------------------------------
// 1.登录

type LoginRequest struct {
	Username string
	Password string
}

// 登录成功后需要 Set-Cookie("camp-session", ${value})
// 密码错误范围密码错误状态码

type LoginResponse struct {
	Code ErrNo
	Data struct {
		UserID string
	}
}

// 2.登出

type LogoutRequest struct{}

// 登出成功需要删除 Cookie

type LogoutResponse struct {
	Code ErrNo
}

// 3.WhoAmI 接口，用来测试是否登录成功，只有此接口需要带上 Cookie

type WhoAmIRequest struct {
}

// 用户未登录请返回用户未登录状态码

type WhoAmIResponse struct {
	Code ErrNo
	Data TMember
}

// -------------------------------------
// 图片信息

// 1.获取成员信息

type GetDataRequest struct {
	DataID string
}

// 如果用户已删除请返回已删除状态码，不存在请返回不存在状态码

type GetDataResponse struct {
	Code ErrNo
	Data TData
}

// 2.批量获取成员信息

type GetDataListRequest struct {
	Offset int
	Limit  int
}

type GetDataListResponse struct {
	Code ErrNo
	Data struct {
		MemberList []TData
	}
}
