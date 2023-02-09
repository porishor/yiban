package app

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var search_mod bool

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "展示所有用户",
	Long:  "展示所有用户，只有手机号",
	Run: func(cmd *cobra.Command, args []string) {
		if search_mod {
			for _, user := range users {
				if len(args) > 0 && user.Username == args[0] {
					log.Println(color.HiBlueString("找到了用户: "), args[0])
					return
				}
			}
			log.Println(color.HiRedString("没有找到用户: "), args[0])
		} else {
			fmt.Printf(color.HiGreenString("%8s\t%15s\n", "序号", "手机号"))
			for index, user := range users {
				fmt.Printf(color.HiYellowString("%8d\t%15s\n", index+1, user.Username))
			}
		}
	},
}
