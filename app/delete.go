package app

import (
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var cookie bool

var deleteCmd = &cobra.Command{
	Use:   "del",
	Short: "del 用户名,根据用户名删除指定用户",
	Long: `del 根据用户名删除指定用户,
		 -cookie 删除指定用户的cookie`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalln(color.HiRedString("参数数量错误，只需要一个参数（用户名）"))
		}
		username := args[0]
		for index, user := range users {
			if user.Username == username {
				// 存在存在cookie flag
				if cookie {
					user.Cookie.CSRF = ""
					user.Cookie.PHPSESSID = ""
				} else {
					users = append(users[:index], users[index+1:]...)
				}
				err := WriteUsers(users)
				if err != nil {
					log.Fatalln(color.HiRedString("删除目标失败"))
				}
				log.Println(color.GreenString("发现了目标"))
			}
		}
	},
}
