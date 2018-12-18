package flowable

import (
	"encoding/base64"
	"os"

	"github.com/imroc/req"
	"github.com/mitchellh/mapstructure"
)

// Service task service definition
type Service interface {
	GetTaskForm(string, string) (map[string]interface{}, error)
	StartProcess(StartProcessForm) (string, error)
	GetUserTasks(string, TaskListQuery) (*FTaskResult, error)
	GetProcessTasks(TaskListQuery) (*FTaskResult, error)
	GetUserProcesses(string, ProcessListQuery) (*FProcessResult, error)
	GetProcess(string) (*FProcess, error)
	SubmitTask(SubmitTaskForm) error
	SubmitTaskAction(string, SubmitTaskActionForm) (*FTask, error)
	CreateUser(NewUserForm) (*UserInfo, error)
	GetUsers() ([]*UserInfo, error)
	CreateAttachment(string, string, string, string, *os.File) (*Attachment, error)
	GetAttachmentFromTask(string) (*AttachmentPaginate, error)
}

type service struct {
	Config Config
}

// NewService 新建任务服务
func NewService(conf Config) Service {
	return &service{Config: conf}
}

func (s *service) GetTaskForm(taskID string, formDefinitionKey string) (map[string]interface{}, error) {
	sEnc := base64.StdEncoding.EncodeToString([]byte(s.Config.RestAccount + ":" + s.Config.RestPasswd))
	header := req.Header{
		"Accept":        "application/json",
		"Authorization": "Basic " + sEnc,
	}
	url := s.Config.Addr + "flowable-task/form-api/form/form-instance-model"
	body := req.BodyJSON(map[string]string{
		"taskId":            taskID,
		"formDefinitionKey": formDefinitionKey,
	})
	r, err := req.Post(url, header, body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	r.ToJSON(&result)
	return result, nil
}

// 创建一个流程，传入流程定义id和表单数据，返回流程的id, processID
// https://www.flowable.org/docs/userguide/index.html#_start_a_process_instance_2
func (s *service) StartProcess(form StartProcessForm) (string, error) {
	sEnc := base64.StdEncoding.EncodeToString([]byte(s.Config.RestAccount + ":" + s.Config.RestPasswd))
	header := req.Header{
		"Accept":        "application/json",
		"Authorization": "Basic " + sEnc,
	}
	url := s.Config.Addr + "flowable-task/process-api/runtime/process-instances/"
	r, err := req.Post(url, header, req.BodyJSON(&form))
	if err != nil {
		return "", err
	}
	var result map[string]interface{}
	r.ToJSON(&result)
	return result["id"].(string), nil
}

// GetUserTasks 取指定用户名下任务列表
// https://www.flowable.org/docs/userguide/index.html#restHistoricTaskInstancesGet
func (s *service) GetUserTasks(state string, query TaskListQuery) (*FTaskResult, error) {
	sEnc := base64.StdEncoding.EncodeToString([]byte(s.Config.RestAccount + ":" + s.Config.RestPasswd))
	header := req.Header{
		"Accept":        "application/json",
		"Authorization": "Basic " + sEnc,
	}
	var params req.Param
	mapstructure.Decode(query, &params)
	if state == "open" {
		params["finished"] = false
	} else if state == "completed" {
		params["finished"] = true
	}
	url := s.Config.Addr + "flowable-task/process-api/history/historic-task-instances/"
	r, err := req.Get(url, header, params)
	if err != nil {
		return nil, err
	}
	result := new(FTaskResult)
	r.ToJSON(result)
	// 把assignee替换为用户对象
	var userIds []string
	for _, row := range result.Data {
		if row.Assignee != "" {
			userIds = append(userIds, row.Assignee)
		}
	}
	users, _ := s.GetUsersByIDs(userIds)
	for i, row := range result.Data {
		for _, user := range users {
			if row.Assignee == user.ID {
				result.Data[i].AssigneeUser = user
			}
		}
	}
	return result, nil
}

// GetProcessTasks 取指定流程下的任务列表
func (s *service) GetProcessTasks(query TaskListQuery) (*FTaskResult, error) {
	sEnc := base64.StdEncoding.EncodeToString([]byte(s.Config.RestAccount + ":" + s.Config.RestPasswd))
	header := req.Header{
		"Accept":        "application/json",
		"Authorization": "Basic " + sEnc,
	}
	var params req.Param
	mapstructure.Decode(query, &params)
	url := s.Config.Addr + "flowable-task/process-api/history/historic-task-instances/"
	r, err := req.Get(url, header, params)
	if err != nil {
		return nil, err
	}
	result := new(FTaskResult)
	r.ToJSON(result)
	// 把assignee替换为用户对象
	var userIds []string
	for _, row := range result.Data {
		if row.Assignee != "" {
			userIds = append(userIds, row.Assignee)
		}
	}
	users, _ := s.GetUsersByIDs(userIds)
	for i, row := range result.Data {
		for _, user := range users {
			if row.Assignee == user.ID {
				result.Data[i].AssigneeUser = user
			}
		}
	}
	return result, nil
}

// GetUserProcesses 获取用户名下的流程
func (s *service) GetUserProcesses(state string, query ProcessListQuery) (*FProcessResult, error) {
	sEnc := base64.StdEncoding.EncodeToString([]byte(s.Config.RestAccount + ":" + s.Config.RestPasswd))
	header := req.Header{
		"Accept":        "application/json",
		"Authorization": "Basic " + sEnc,
	}
	var params req.Param
	mapstructure.Decode(query, &params)
	// state=all，显示所有
	if state == "running" {
		params["finished"] = false
	} else if state == "completed" {
		params["finished"] = true
	}
	url := s.Config.Addr + "flowable-task/process-api/history/historic-process-instances/"
	r, err := req.Get(url, header, params)
	if err != nil {
		return nil, err
	}
	result := new(FProcessResult)
	r.ToJSON(result)
	var userIds []string
	for _, row := range result.Data {
		userIds = append(userIds, row.StartUserID)
	}
	users, _ := s.GetUsersByIDs(userIds)
	for i, row := range result.Data {
		for _, user := range users {
			if row.StartUserID == user.ID {
				result.Data[i].StartedBy = user
			}
		}
	}
	return result, nil
}

// SubmitTask 提交任务关联的表单
func (s *service) SubmitTask(form SubmitTaskForm) error {
	sEnc := base64.StdEncoding.EncodeToString([]byte(s.Config.RestAccount + ":" + s.Config.RestPasswd))
	header := req.Header{
		"Accept":        "application/json",
		"Authorization": "Basic " + sEnc,
	}
	url := s.Config.Addr + "flowable-task/process-api/form/form-data/"
	r, err := req.Post(url, header, req.BodyJSON(form))
	if err != nil {
		return err
	}
	var result map[string]interface{}
	r.ToJSON(&result)
	return nil
}

// SubmitTaskAction 提交任务关联的表单
// https://www.flowable.org/docs/userguide/index.html#_task_actions
func (s *service) SubmitTaskAction(taskID string, form SubmitTaskActionForm) (*FTask, error) {
	sEnc := base64.StdEncoding.EncodeToString([]byte(s.Config.RestAccount + ":" + s.Config.RestPasswd))
	header := req.Header{
		"Accept":        "application/json",
		"Authorization": "Basic " + sEnc,
	}
	url := s.Config.Addr + "flowable-task/process-api/runtime/tasks/" + taskID
	r, err := req.Post(url, header, req.BodyJSON(form))
	if err != nil {
		return nil, err
	}
	var result *FTask
	r.ToJSON(result)
	return result, nil
}

func (s *service) CreateUser(newUser NewUserForm) (*UserInfo, error) {
	sEnc := base64.StdEncoding.EncodeToString([]byte(s.Config.RestAccount + ":" + s.Config.RestPasswd))
	header := req.Header{
		"Accept":        "application/json",
		"Authorization": "Basic " + sEnc,
	}
	url := s.Config.Addr + "flowable-task/process-api/identity/users/"
	r, err := req.Post(url, header, req.BodyJSON(newUser))
	if err != nil {
		return nil, err
	}
	var result *UserInfo
	r.ToJSON(result)
	return result, nil
}

func (s *service) CreateAttachment(taskID string, filedName string, fileName string, mimeType string, content *os.File) (*Attachment, error) {
	sEnc := base64.StdEncoding.EncodeToString([]byte(s.Config.RestAccount + ":" + s.Config.RestPasswd))
	header := req.Header{
		"Accept":        "application/json",
		"Authorization": "Basic " + sEnc,
	}
	url := s.Config.Addr + "flowable-task/content-api/content-service/content-items/"
	param := req.Param{
		"taskId":   taskID,
		"name":     filedName,
		"mimeType": mimeType,
	}
	r, err := req.Post(url, header, param, req.FileUpload{
		File:      content,
		FieldName: filedName,
		FileName:  fileName,
	})
	if err != nil {
		return nil, err
	}
	var result *Attachment
	r.ToJSON(&result)
	return result, nil
}

func (s *service) GetAttachmentFromTask(taskID string) (*AttachmentPaginate, error) {
	sEnc := base64.StdEncoding.EncodeToString([]byte(s.Config.RestAccount + ":" + s.Config.RestPasswd))
	header := req.Header{
		"Accept":        "application/json",
		"Authorization": "Basic " + sEnc,
	}
	url := s.Config.Addr + "flowable-task/content-api/content-service/content-items/"
	param := req.Param{
		"taskId": taskID,
	}
	r, err := req.Get(url, header, param)
	if err != nil {
		return nil, err
	}
	var result *AttachmentPaginate
	r.ToJSON(&result)
	return result, nil
}

func (s *service) GetProcess(procID string) (*FProcess, error) {
	sEnc := base64.StdEncoding.EncodeToString([]byte(s.Config.RestAccount + ":" + s.Config.RestPasswd))
	header := req.Header{
		"Accept":        "application/json",
		"Authorization": "Basic " + sEnc,
	}
	url := s.Config.Addr + "flowable-task/process-api/history/historic-process-instances/" + procID
	r, err := req.Get(url, header)
	if err != nil {
		return nil, err
	}
	var result *FProcess
	r.ToJSON(&result)
	startUserID := result.StartUserID
	users, err := s.GetUsersByIDs([]string{startUserID})
	if err == nil {
		result.StartedBy = users[0]
	}
	return result, nil
}

func (s *service) GetUsers() ([]*UserInfo, error) {
	// 从flowable中获取用户列表
	sEnc := base64.StdEncoding.EncodeToString([]byte(s.Config.RestAccount + ":" + s.Config.RestPasswd))
	header := req.Header{
		"Accept":        "application/json",
		"Authorization": "Basic " + sEnc,
	}
	url := s.Config.Addr + "flowable-task/process-api/identity/users/"
	r, err := req.Get(url, header)
	if err != nil {
		return nil, err
	}
	var paginateResult map[string]interface{}
	r.ToJSON(&paginateResult)
	data := paginateResult["data"]
	var result []*UserInfo
	mapstructure.Decode(data, &result)
	return result, nil
}

func (s *service) GetUsersByIDs(ids []string) ([]*UserInfo, error) {
	var result []*UserInfo
	users, err := s.GetUsers()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		for _, id := range ids {
			if user.ID == id {
				result = append(result, user)
			}
		}
	}
	return result, nil
}
