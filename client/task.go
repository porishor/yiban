package client

// API:https://api.uyiban.com/officeTask/client/index/uncompletedList 查询未完成的任务，返回[]Task
type Task struct {
	TaskID       string `json:"TaskId"`
	OrgID        string `json:"OrgId"`
	TimeoutState int64  `json:"TimeoutState"`
	State        int64  `json:"State"`
	Title        string `json:"Title"`
	Type         int64  `json:"Type"`
	StartTime    int64  `json:"StartTime"`
	EndTime      int64  `json:"EndTime"`
}

// API: https://api.uyiban.com/officeTask/client/index/detail?TaskId={}&CSRF={}
type TaskDetail struct {
	ID                 string        `json:"Id"`
	Title              string        `json:"Title"`
	Type               int64         `json:"Type"`
	ReceiverSendBack   int64         `json:"ReceiverSendBack"`
	TypeID             string        `json:"TypeId"`
	AllowSendBackHours int64         `json:"AllowSendBackHours"`
	AllowSendBackTimes int64         `json:"AllowSendBackTimes"`
	PubOrgName         string        `json:"PubOrgName"`
	PubPersonName      string        `json:"PubPersonName"`
	Content            string        `json:"Content"`
	AllowTimeout       int64         `json:"AllowTimeout"`
	AttachmentIDS      string        `json:"AttachmentIds"`
	PubOrgID           string        `json:"PubOrgId"`
	IsPubPersonShow    int64         `json:"IsPubPersonShow"`
	CreateTime         int64         `json:"CreateTime"`
	StartTime          int64         `json:"StartTime"`
	EndTime            int64         `json:"EndTime"`
	Displayed          int64         `json:"Displayed"`
	TimeState          int64         `json:"TimeState"`
	WFID               string        `json:"WFId"`
	AttachmentList     []interface{} `json:"AttachmentList"`
	WorkflowState      int64         `json:"WorkflowState"`
	InitiateID         string        `json:"InitiateId"`
	RetreatReason      string        `json:"RetreatReason"`
	EvaluationState    int64         `json:"EvaluationState"`
	EvaluationReason   string        `json:"EvaluationReason"`
	State              interface{}   `json:"State"`
	RecordState        int64         `json:"RecordState"`
	CanSendBack        int64         `json:"CanSendBack"`
	SendBackEndTime    int64         `json:"SendBackEndTime"`
	IsLost             int64         `json:"IsLost"`
	Feedback           int64         `json:"Feedback"`
	ChangeReason       string        `json:"ChangeReason"`
}

// ProcessDetail API:https://api.uyiban.com/workFlow/c/my/getProcessDetail?WFId={}&CSRF={}
type ProcessDetail struct {
	List       []List     `json:"list"`
	Kv         Kv         `json:"kv"`
	PersonInfo PersonInfo `json:"personInfo"`
}

type Kv struct {
	Person interface{}     `json:"person"`
	Org    interface{}     `json:"org"`
	Role   map[string]Role `json:"role"`
}

type Role struct {
	ID       string `json:"Id"`
	RoleName string `json:"RoleName"`
}

type List struct {
	ID             string        `json:"Id"`
	UniversityID   string        `json:"UniversityId"`
	WFID           string        `json:"WFId"`
	PIndex         int64         `json:"PIndex"`
	Cond           interface{}   `json:"Cond"`
	RulePersonCond string        `json:"RulePersonCond"`
	Flow           []interface{} `json:"Flow"`
	CCTrigger      string        `json:"CCTrigger"`
	Cc             Cc            `json:"CC"`
	State          int64         `json:"State"`
	CreateTime     int64         `json:"CreateTime"`
}

type Cc struct {
	Rulerole []string `json:"rulerole"`
}

type PersonInfo struct {
	ID             string      `json:"Id"`
	UniversityID   string      `json:"UniversityId"`
	Name           string      `json:"Name"`
	Gender         int64       `json:"Gender"`
	CollegeID      string      `json:"CollegeId"`
	FacultyID      interface{} `json:"FacultyId"`
	ProfessionID   string      `json:"ProfessionId"`
	SpecialtyID    interface{} `json:"SpecialtyId"`
	ClassID        string      `json:"ClassId"`
	Grade          int64       `json:"Grade"`
	EducationLevel int64       `json:"EducationLevel"`
	StudyYear      string      `json:"StudyYear"`
	Number         string      `json:"Number"`
	ExamNo         interface{} `json:"ExamNo"`
	AdmissionNo    interface{} `json:"AdmissionNo"`
	IDCardNo       interface{} `json:"IDCardNo"`
	Mobile         interface{} `json:"Mobile"`
	MobileState    int64       `json:"MobileState"`
	Email          interface{} `json:"Email"`
	EmailState     int64       `json:"EmailState"`
	State          int64       `json:"State"`
	IsFreeze       int64       `json:"IsFreeze"`
	Campus         string      `json:"Campus"`
	IsArchived     int64       `json:"IsArchived"`
	Remark         interface{} `json:"Remark"`
	CreateTime     int64       `json:"CreateTime"`
	UpdateTime     int64       `json:"UpdateTime"`
	PersonType     string      `json:"PersonType"`
	CollegeName    string      `json:"CollegeName"`
	FacultyName    string      `json:"FacultyName"`
	ProfessionName string      `json:"ProfessionName"`
	SpecialtyName  string      `json:"SpecialtyName"`
	ClassName      string      `json:"ClassName"`
}

// 上次所填内容, API:https://api.uyiban.com/workFlow/c/work/getLastInitiate?WFId={}&CSRF={}

type LastFormContent struct {
	WFName           string        `json:"WFName"`
	WFOtherConfig    []interface{} `json:"WFOtherConfig"`
	Cc2              []Cc2         `json:"CC"`
	Process          Process       `json:"Process"`
	Approved         []interface{} `json:"Approved"`
	HasEditionAccess bool          `json:"HasEditionAccess"`
	HasInitiateLog   bool          `json:"HasInitiateLog"`
	Initiate         Initiate      `json:"Initiate"`
	InitiateExtend   []interface{} `json:"InitiateExtend"`
	ReuseFillState   int64         `json:"ReuseFillState"`
}

type Cc2 struct {
	CCPersonID    string `json:"CCPersonId"`
	CCPersonLabel string `json:"CCPersonLabel"`
}

type Initiate struct {
	ID             string         `json:"Id"`
	SerialNo       string         `json:"SerialNo"`
	UniversityID   string         `json:"UniversityId"`
	PersonID       string         `json:"PersonId"`
	PersonInfo2    PersonInfo2    `json:"PersonInfo"`
	WFID           string         `json:"WFId"`
	ProcessID      string         `json:"ProcessId"`
	FormDataJSON   []FormDataJSON `json:"FormDataJson"`
	ExtendDataJSON ExtendDataJSON `json:"ExtendDataJson"`
	WorkNode       int64          `json:"WorkNode"`
	State          int64          `json:"State"`
	StateTime      string         `json:"StateTime"`
	CreateTime     int64          `json:"CreateTime"`
}

type ExtendDataJSON struct {
	TaskID  string    `json:"TaskId"`
	Title   string    `json:"title"`
	Content []Content `json:"content"`
}

type Content struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type FormDataJSON struct {
	ID        string      `json:"id"`
	Label     string      `json:"label"`
	Value     interface{} `json:"value"`
	Component string      `json:"component"`
}

type ValueClass struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Address  string `json:"address"`
}

type PersonInfo2 struct {
	College        string      `json:"College"`
	Profession     string      `json:"Profession"`
	Class          string      `json:"Class"`
	Grade          int64       `json:"Grade"`
	Campus         string      `json:"Campus"`
	EducationLevel int64       `json:"EducationLevel"`
	StudyYear      string      `json:"StudyYear"`
	Number         string      `json:"Number"`
	Name           string      `json:"Name"`
	Gender         int64       `json:"Gender"`
	PersonType     string      `json:"PersonType"`
	Mobile         interface{} `json:"Mobile"`
}

type Process struct {
	WFID      string        `json:"WFId"`
	Flow      []interface{} `json:"Flow"`
	CCTrigger string        `json:"CCTrigger"`
}

type ValueUnion struct {
	String      *string
	StringArray []string
	ValueClass  *ValueClass
}
