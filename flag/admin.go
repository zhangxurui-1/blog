package flag

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"golang.org/x/term"
	"os"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/utils"
	"syscall"
)

// Admin 创建一个管理员用户
func Admin() error {
	var user database.User
	// 输入邮箱
	fmt.Print("Enter email:")
	var email string
	_, err := fmt.Scanln(&email)
	if err != nil {
		return fmt.Errorf("failed to read email: %v", err)
	}
	user.Email = email

	// 获取标准输入的文件描述符
	fd := int(syscall.Stdin)
	// 关闭回显，使密码不会在终端显示
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	// 恢复终端状态
	defer term.Restore(fd, oldState)

	// 输入密码
	fmt.Print("Enter password:")
	password, err := readPassword()
	fmt.Println()
	if err != nil {
		return err
	}

	fmt.Print("Confirm password:")
	rePassword, err := readPassword()
	fmt.Println()
	if err != nil {
		return err
	}

	// 检查两次密码输入是否一致
	if password != rePassword {
		return errors.New("passwords do not match")
	}

	// 检查密码长度是否符合要求
	if len(password) < 8 || len(rePassword) > 20 {
		return errors.New("password length must be between 8 and 20 characters")
	}

	// 填充用户数据
	user.UUID = uuid.Must(uuid.NewV4())
	user.Username = global.Config.Website.Name
	user.Password = utils.BcryptHash(password)
	user.RoleID = appTypes.Admin
	user.Avatar = "/image/avatar.jpg"
	user.Address = global.Config.Website.Address

	// 在数据库中创建管理员用户
	if err := global.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

// readPassword 读取密码并避免回显
func readPassword() (string, error) {
	var password string
	var buf [1]byte

	// 持续读取字符直到遇到换行符为止
	for {
		_, err := os.Stdin.Read(buf[:])
		if err != nil {
			return "", err
		}

		char := buf[0]
		if char == '\n' || char == '\r' {
			break
		}
		password = password + string(char)
	}
	return password, nil
}
