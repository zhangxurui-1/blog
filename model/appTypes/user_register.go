package appTypes

import "encoding/json"

// Register 用户注册来源
type Register int

const (
	Email Register = iota
	QQ
	Beta
)

// MarshalJSON 实现了 json.Marshaler 接口
func (r Register) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

// UnmarshalJSON 实现了 json.Unmarshaler 接口
func (r *Register) UnmarshalJSON(bytes []byte) error {
	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}
	*r = ToRegister(s)
	return nil
}

// 返回 Register 的字符串表示
func (r Register) String() string {
	var str string
	switch r {
	case Email:
		str = "邮箱"
	case QQ:
		str = "QQ"
	case Beta:
		str = "Beta"
	default:
		str = "未知"
	}

	return str
}

// ToRegister 将字符串转为 Register
func ToRegister(str string) Register {
	switch str {
	case "邮箱":
		return Email
	case "QQ":

		return QQ
	case "Beta":
		return Beta
	default:
		return -1

	}
}
