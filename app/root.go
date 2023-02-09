package app

import (
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var users []*User

func App() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "yiban",
		Short: "填写易班日报",
		Long: `填写易班日报，支持多个用户。会将昨天填写的易班日报内容填写到今天
		使用方法: yiban run 或者 yiban 会直接使用$APPDADA/yiban/user.csv中的用户
		         yiban add 用户名 密码 添加用户
				 yiban del 用户名  删除用户
				 yiban show 查看所有用户`,
		Run: func(cmd *cobra.Command, args []string) {
			RunSubmit()
		},
	}

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "运行填写日报的程序",
		Long:  "找到$APPDATA/yiban/user.csv，运行程序",
		Run: func(cmd *cobra.Command, args []string) {
			RunSubmit()
		},
	})
	return rootCmd
}
func init() {
	log.Println(color.BlueString("初始化....."))
	var err error
	users, err = ReadConf()
	if err != nil {
		log.Fatalln(color.HiRedString("读取配置失败"))
	}
	deleteCmd.Flags().BoolVarP(&cookie, "cookie", "c", false, "是否删除cookie")
	showCmd.Flags().BoolVarP(&search_mod, "search", "s", false, "查找用户是否存在")
}
