package app

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"yiban/client"

	"github.com/fatih/color"
)

type UserCookie struct {
	CSRF      string `json:"csrfToken"`
	PHPSESSID string `json:"phpsessid"`
}
type User struct {
	Username string     `json:"username"`
	Password string     `json:"password"`
	Cookie   UserCookie `json:"cookie"`
	// 列举那些任务名称
	Include string `json:"include"`
	Exclue  string `json:"exclue"`
}

func conf_path() string {
	appdat := os.Getenv("APPDATA")
	return path.Join(appdat, "yiban\\user.csv")
}
func ReadConf() ([]*User, error) {

	opencast, err := os.Open(conf_path())
	if err != nil {
		return nil, err
	}
	defer opencast.Close()
	csvReader := csv.NewReader(opencast)
	read, err := csvReader.Read()
	if len(read) != 6 {
		log.Printf("文件格式错误:%s", err)
		return nil, err
	}
	// log.Printf("\n%10s\t%10s\t%32s\t%32s\n", read[0], read[1], read[2], read[3])
	readAll, err := csvReader.ReadAll()
	users := make([]*User, len(readAll))
	for index, line := range readAll {
		if len(line) != 6 {
			log.Printf("文件格式错误:%s，位于第%d行", err, index+1)
			continue
		}
		user := &User{
			Username: line[0],
			Password: line[1],
			Cookie: UserCookie{
				CSRF:      line[2],
				PHPSESSID: line[3],
			},
			Include: line[4],
			Exclue:  line[5],
		}
		users[index] = user
	}
	return users, nil
}

func WriteUsers(users []*User) error {
	//OpenFile读取文件，不存在时则创建，使用追加模式
	File, err := os.OpenFile(conf_path(), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln("文件打开失败！")
	}
	defer File.Close()

	//创建写入接口
	WriterCsv := csv.NewWriter(File)

	WriterCsv.Write([]string{"用户名", "密码", "cookie(csrf_token)", "cookie(phpseeid)", "包含的关键词(用、分割)", "排除的关键词"})

	for index, user := range users {
		// log.Printf("写入第%d条数据\n", index)
		str := []string{user.Username, user.Password, user.Cookie.CSRF, user.Cookie.PHPSESSID, user.Include, user.Exclue} //需要写入csv的数据，切片类型
		//写入一条数据，传入数据为切片(追加模式)
		err1 := WriterCsv.Write(str)
		if err1 != nil {
			log.Printf("WriterCsv写入文件失败,位于第%d行\n", index)
			return err1
		}
	}
	WriterCsv.Flush()
	return nil
}
// 读取配置文件并且提交全部日报
func RunSubmit() {

	// file_name := fmt.Sprintf(time.Now().Format("2006年01月02日易班日志.log"))
	// file, err := os.OpenFile(path.Join("D:\\workspace\\yiban\\log\\", file_name), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	//     log.Fatalln("Faild to open error logger file:", err)
	// }
	// log.SetOutput(file)
	users, err1 := ReadConf()
	if err1 != nil {
		log.Println(color.RedString("读取文件失败，是否创建新文件.回复Y(Yes) or N(No)?"))
		var option byte
		fmt.Scanf("%c", &option)
		if option == 'Y' {
			origin := &User{
				Username: "***",
				Password: "***",
				Cookie: UserCookie{
					CSRF:      "",
					PHPSESSID: "",
				},
				Include: "日报",
				Exclue:  "",
			}
			err2 := WriteUsers([]*User{
				origin, origin, origin,
			})
			if err2 != nil {
				log.Printf("创建新文件失败了!\n")
			} else {
				log.Printf("成功\n")
			}
		}
	}
	for index, user := range users {
		if user.Username == "" || user.Password == "" {
			log.Printf("用户名或者密码不能为空，位于第%d行\n", index+1)
			continue
		}
		if user.Cookie.CSRF != "" && user.Cookie.PHPSESSID != "" {
			helper, err := client.CreateHelperWithCookie([]*http.Cookie{
				{
					Name:  "csrf_token",
					Value: user.Cookie.CSRF,
				},
				{
					Name:  "PHPSESSID",
					Value: user.Cookie.PHPSESSID,
				},
			})
			if err != nil || helper.Name == "" {
				log.Printf("通过cookie创建失败，将通过用户名密码尝试\n")
				helper, err = client.CreateHelperWithPassword(strings.TrimSpace(user.Username),
					strings.TrimSpace(user.Password))

				if err != nil {
					log.Println(color.RedString("登陆失败"), "进行最后一次尝试:\n", err)
					helper, err = client.CreateHelperWithPassword(user.Username, user.Password)
					if err != nil {
						log.Println(color.RedString("登陆失败"), "请检查用户名和密码:", user.Username, user.Password)
					} else {
						log.Printf(color.GreenString("登陆成功，姓名:%s，大学:%s\n", helper.Name, helper.University))
					}
				} else {
					log.Printf(color.GreenString("登陆成功，姓名:%s，大学:%s\n", helper.Name, helper.University))
				}
			} else {
				log.Printf(color.GreenString("登陆成功，姓名:%s，大学:%s\n", helper.Name, helper.University))
			}
			user.Cookie.CSRF = *client.ExtractCSRF(helper.Cookies)
			user.Cookie.PHPSESSID = *client.ExactCookie(helper.Cookies, "PHPSESSID")
			helper.SubmitAll(user.Include, user.Exclue)
		} else {
			helper, err := client.CreateHelperWithPassword(user.Username, user.Password)

			if err != nil {
				log.Println(color.RedString("登陆失败"), "进行最后一次尝试:\n", err)
				helper, err = client.CreateHelperWithPassword(user.Username, user.Password)
				if err != nil {
					log.Println(color.RedString("登陆失败"), "请检查用户名和密码:", user.Username, user.Password)
					continue
				}
			}
			user.Cookie.CSRF = *client.ExtractCSRF(helper.Cookies)
			user.Cookie.PHPSESSID = *client.ExactCookie(helper.Cookies, "PHPSESSID")
			helper.SubmitAll(user.Include, user.Exclue)
		}
		log.Println(color.HiBlueString("*******************分割线*********************"))
	}
	err := WriteUsers(users)
	if err != nil {
		log.Printf(color.RedString("写入失败\n"))
		log.Printf("\n%10s\t%10s\t%32s\t%32s\t%20s\t%20s\n", "用户名", "密码", "cookie(csrf_token)", "cookie(phpseeid)", "包含的关键词(用,分割)", "排除的关键词")
		for _, user := range users {
			log.Printf("%10s\t%10s\t%32s\t%32s\t%20s\t%20s\n", user.Username, user.Password, user.Cookie.CSRF, user.Cookie.PHPSESSID, user.Include, user.Exclue)
		}
	}
}
