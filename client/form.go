package client

import (
	"encoding/json"
	"fmt"
)

type Form struct {
	WFID          string `json:"WFId"`
	Data          string `json:"Data"`
	WfprocessID   string `json:"WfprocessId"`
	Extend        string `json:"Extend"`
	CustomProcess string `json:"CustomProcess"`
}

func genFormFromLastData(initial *Initiate, WFId string, WfprocessId string,
	taskId string, taskTitle string) (*Form, error) {
	forms := initial.FormDataJSON
	data := make(map[string]interface{})

	for _, form := range forms {
		if form.Value != nil {
			data[form.ID] = form.Value
			// fmt.Printf("\nId:%s, label:%s, value:%v\n", form.ID, form.Label, form.Value)
		}
	}
	extend := fmt.Sprintf(`{"TaskId":"%s","title":"任务信息","content":[{"label":"任务名称","value":"%s"},{"label":"发布机构","value":"学生工作处"}]}`, taskId, taskTitle)
	customProcess := `{"ApplyPersonIds":[],"CCPersonId":[]}`
	data_str, err := json.Marshal(&data)
	if err != nil {
		return nil, fmt.Errorf("生成Form时序列化错误：%s", err)
	}
	return &Form{
		WFID:          WFId,
		Data:          string(data_str),
		WfprocessID:   WfprocessId,
		Extend:        extend,
		CustomProcess: customProcess,
	}, nil

}
