package flowable

// StartProcessForm 开始流程所需要的表单定义
type StartProcessForm struct {
	ProcessDefinitionID string `json:"processDefinitionId"`
}

// FormVariable 单个表单一行提交时的内容
type FormVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// FTaskResult 一个任务的表示
type FTaskResult struct {
	Data  []FTask `json:"data"`
	Total int     `json:"total"`
	Start int     `json:"start"`
	Sort  string  `json:"sort"`
	Order string  `json:"order"`
	Size  int     `json:"size"`
}

// FTask flowable 任务对象
type FTask struct {
	ID                  string         `json:"id"`
	Name                string         `json:"name"`
	Assignee            string         `json:"assignee"`
	AssigneeUser        *UserInfo      `json:"assignee_user"`
	ForkKey             string         `json:"formKey"`
	CreateTime          string         `json:"createTime"`
	TaskDefinitionKey   string         `json:"taskDefinitionKey"`
	ExecutionID         string         `json:"executionId"`
	ProcessInstanceID   string         `json:"processInstanceId"`
	ProcessDefinitionID string         `json:"processDefinitionId"`
	Description         string         `json:"description"`
	StartTime           string         `json:"startTime"`
	EndTime             string         `json:"endTime"`
	ClaimTime           string         `json:"claimTime"`
	DueDate             string         `json:"dueDate"`
	Variables           []FormVariable `json:"variables"`
}

// FProcess 流程对象
type FProcess struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	URL                  string    `json:"url"`
	BusinessKey          string    `json:"businessKey"`
	Suspended            string    `json:"suspended"`
	ProcessDefinitionURL string    `json:"processDefinitionUrl"`
	ActivityID           string    `json:"activityId"`
	StartTime            string    `json:"startTime"`
	EndTime              string    `json:"endTime"` // 可通过此字段判断流程是否结束
	StartUserID          string    `json:"startUserId"`
	StartedBy            *UserInfo `json:"startedBy"` // 后加的字段，默认接口没返回
}

// FProcessResult 流程列表查询结果
type FProcessResult struct {
	Data  []FProcess `json:"data"`
	Total int        `json:"total"`
	Start int        `json:"start"`
	Sort  string     `json:"sort"`
	Order string     `json:"order"`
	Size  int        `json:"size"`
}

// SubmitTaskForm 提交任务表单时传的参数
type SubmitTaskForm struct {
	TaskID     string         `json:"task_id"`
	Properties []FormVariable `json:"properties"`
}

// SubmitTaskActionForm 提交任务动作表单
type SubmitTaskActionForm struct {
	Action    string         `json:"action"`
	Variables []FormVariable `json:"variables"`
}

// NewUserForm 新建一个用户到flowable所填写的字段
type NewUserForm struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// UserInfo 单个用户的信息
type UserInfo struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	URL       string `json:"url"`
	Email     string `json:"email"`
}

// Attachment 附件上传后的结果
type Attachment struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	MimeType          string `json:"mimeType,omitempty"`
	TaskID            string `json:"taskId"`
	ProcessInstanceID string `json:"processInstanceId,omitempty"`
	ContentStoreID    string `json:"contentStoreId"`
	ContentStoreName  string `json:"contentStoreName"`
	ContentAvailable  bool   `json:"contentAvailable"`
	Created           string `json:"created"`
	CreatedBy         string `json:"createdBy,omitempty"`
	LastModified      string `json:"lastModified"`
	LastModifiedBy    string `json:"lastModifiedBy,omitempty"`
	URL               string `json:"url"`
}

// AttachmentPaginate 附件分页列表
type AttachmentPaginate struct {
	Data  []Attachment `json:"data"`
	Total int          `json:"total"`
	Start int          `json:"start"`
	Sort  string       `json:"sort"`
	Order string       `json:"order"`
	Size  int          `json:"size"`
}

type TaskListQuery struct {
	ProcessInstanceID string `json:"processInstanceID,omitempty"`
	TaskAssignee      string `json:"taskAssignee,omitempty"`
	Start             int    `json:"start"`
	Size              int    `json:"size"`
}

type ProcessListQuery struct {
	InvolvedUser string `json:"involvedUser"`
	Start        int    `json:"start"`
	Size         int    `json:"size"`
}
