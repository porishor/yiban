package app

import (
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "添加用户",
	Long:  "添加用户 yiban add 手机号(用户名) 密码",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Fatalln("缺乏参数,需要 add 用户名 密码")
			return
		}
		usernmae := args[0]
		password := args[1]
		
		users = append(users, &User{
			Username: usernmae,
			Password: password,
		})
		err := WriteUsers(users)
		if err != nil {
			log.Println(color.HiRedString("写入数据失败"))
		}
		log.Println(color.HiGreenString("写入数据成功"))
	},
}
