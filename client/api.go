package client

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"yiban/encrypt"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
)

type Helper struct {
	Cookies       []*http.Cookie
	Name          string
	University    string
	person_id     string
	university_id string
}

func CreateHelperWithCookie(cookies []*http.Cookie) (*Helper, error) {

	csrf_token := ExtractCSRF(cookies)
	phpsession := ExactCookie(cookies, "PHPSESSID")
	if csrf_token == nil || phpsession == nil {
		return nil, fmt.Errorf("\n获取cookies失败:%v\n", cookies)
	}
	req := BuildWithCookie(*csrf_token, *phpsession)
	resp, err := req.Get(AuthURL)
	if err != nil {
		return nil, err
	}
	var resp_json YiBanResponse[UserInfo]
	err = json.Unmarshal(resp.Body(), &resp_json)
	if err != nil {
		return nil, fmt.Errorf("\n(%s)序列化失败,API:%s:%s\n", "用户信息()", AuthURL, err.Error())
	}

	return &Helper{
		Cookies:       cookies,
		Name:          resp_json.Data.PersonName,
		University:    resp_json.Data.UniversityName,
		person_id:     resp_json.Data.PersonID,
		university_id: resp_json.Data.UniversityID,
	}, nil
}
func CreateHelperWithPassword(username string, password string) (*Helper, error) {
	cookies, err := Login(username, password)
	if err != nil {
		return nil, fmt.Errorf("\n登陆失败\n:%s", err.Error())
	}
	return CreateHelperWithCookie(cookies)
}

// 获取从d天之前到现在的未完成的任务
func (helper *Helper) GetUncompletedTasks(d int) ([]Task, error) {
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -d)
	var resp_json YiBanResponse[[]Task]
	helper.get(func(r *resty.Request) {
		r.SetQueryParams(map[string]string{
			"StartTime": fmt.Sprintf("%d-%d-%d 00:00", startTime.Year(), startTime.Month(), startTime.Day()),
			"EndTime":   fmt.Sprintf("%d-%d-%d 23:59", endTime.Year(), endTime.Month(), endTime.Day()),
		})
	}, UnCompletedList, &resp_json)
	return resp_json.Data, nil
}

func (helper *Helper) GetTaskDetail(taskId string) (*TaskDetail, error) {
	var resp_json YiBanResponse[TaskDetail]
	err := helper.get(func(r *resty.Request) {
		r.SetQueryParam("TaskId", taskId)
	}, TaskDetailURL, &resp_json)
	if err != nil {
		return nil, err
	}
	return &resp_json.Data, nil
}

func (helper *Helper) GetProcessDetail(WFId string) (*ProcessDetail, error) {
	var resp_json YiBanResponse[ProcessDetail]
	err := helper.get(func(r *resty.Request) {
		r.SetQueryParam("WFId", WFId)
	}, ProcessDetailURL,
		&resp_json)
	if err != nil {
		return nil, err
	}
	return &resp_json.Data, nil
}

func (helper *Helper) GetLastInitial(WFId string) (*LastFormContent, error) {
	var resp_json YiBanResponse[LastFormContent]
	err := helper.get(func(r *resty.Request) {
		r.SetQueryParam("WFId", WFId)
	}, LastInitial,
		&resp_json)
	if err != nil {
		return nil, err
	}
	return &resp_json.Data, nil
}

func (helper *Helper) get(requestHandler func(*resty.Request), url string, data interface{}) error {
	req := BuildYibanRequest(helper.Cookies)
	if requestHandler != nil {
		requestHandler(req)
	}
	resp, err := req.
		Get(url)
	if err != nil {
		return err
	}
	/* var zh, _ = zhToUnicode(resp.Body())
	fmt.Printf("\n%s\n", zh) */

	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		return fmt.Errorf("序列化失败:\n%s\n", err.Error())
	}
	return nil
}

// 返回提交日报后,服务器的信息
func (helper *Helper) Submit(taskId string) (*YiBanResponse[string], error) {
	taskDetail, err := helper.GetTaskDetail(taskId)
	startTime := time.Unix(taskDetail.StartTime, 0)
	log.Printf("%s\n%s:%s\n%s:%s\n%s:%v\n%s:%s\n", color.BlueString("\n获取任务详细信息成功"),
		color.BlueString("工作流id"), color.HiGreenString(taskDetail.WFID),
		color.BlueString("标题"), color.HiGreenString(taskDetail.Title), color.BlueString("状态"),
		color.HiGreenString("%s", taskDetail.State), color.BlueString("起始时间"),
		color.HiCyanString(startTime.Format("2006-01-02 15:04:05")))
	if time.Now().Before(startTime) {
		return nil, errors.New(color.HiRedString("未到签到时间,TaskId:%s", taskId))
	}
	var wf_id = taskDetail.WFID
	processDetail, err := helper.GetProcessDetail(wf_id)
	if err != nil {
		return nil, err
	}
	var WfprocessID = processDetail.List[0].ID

	last, err := helper.GetLastInitial(wf_id)
	if err != nil {
		return nil, err
	}
	form, err := genFormFromLastData(&last.Initiate, wf_id, WfprocessID, taskId, taskDetail.Title)

	// form_str := fmt.Sprintf(`{"WFId":"%s","Data":"%s","WfprocessId":"%s","Extend":"{\"TaskId\":\"%s\",\"title\":\"任务信息\",\"content\":[{\"label\":\"任务名称\",\"value\":\"%s\"},{\"label\":\"发布机构\",\"value\":\"学生工作处\"}]}","CustomProcess":"{\"ApplyPersonIds\":[],\"CCPersonId\":[]}"}`, wf_id, form.Data, WfprocessID, taskId, taskDetail.Title)
	if err != nil {
		return nil, fmt.Errorf("\n生成form失败;%s\n", err)
	}
	form_byte, err := json.Marshal(form)
	// var form_byte = []byte(form_str)
	if err != nil {
		return nil, fmt.Errorf("\nform的反序列化失败:%s\n", err)
	}
	form1 := encrypt.AesEncryptCBC(form_byte, []byte(encrypt.KEY), []byte(encrypt.IV))
	form2 := base64.StdEncoding.EncodeToString(form1)
	form3 := base64.StdEncoding.EncodeToString([]byte(form2))
	// fmt.Printf("\n加密后结果:%s\n", form3)
	resp, err := BuildYibanRequest(helper.Cookies).
		SetFormData(
			map[string]string{
				"Str": form3,
			}).
		Post("https://api.uyiban.com/workFlow/c/my/apply")

	var msg YiBanResponse[string]
	err = json.Unmarshal(resp.Body(), &msg)
	if err != nil {
		return nil, fmt.Errorf("\n 服务器返回结果json化失败:%s\n", err)
	}
	// fmt.Printf("\ndata:%+v\n",msg.Msg)
	return &msg, nil
}

func (helper *Helper) SubmitAll(include string, exclude string) {
	tasks, err := helper.GetUncompletedTasks(15)
	in_keys := strings.Split(include, "、")
	ex_keys := strings.Split(exclude, "、")
	if err != nil {
		log.Printf("\n获取未完成任务失败:%s\n", err)
		return
	}
	for _, task := range tasks {
		var flag = true
		for _, key := range in_keys {
			if !strings.Contains(task.Title, key) && key != "" {
				flag = false
			}
		}
		for _, key := range ex_keys {
			if strings.Contains(task.Title, key) && key != "" {
				flag = false
			}
		}
		if !flag {
			log.Printf(color.HiYellowString("跳过了:%s\n", task.Title))
			continue
		}
		resp, err := helper.Submit(task.TaskID)
		if err != nil {
			log.Printf("%s(%s) error:%s\n", color.RedString("提交失败"), color.YellowString(task.Title), err)
			continue
		}
		log.Printf("\n提交成功(%s)\n返回结果:%s\n返回消息:%s\n", task.Title, resp.Data, resp.Msg)
	}
}
