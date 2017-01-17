package base

import (
	"encoding/json"
	"fmt"
	md "goERP/models"
	"strconv"
	"strings"
)

// UserController user
type UserController struct {
	BaseController
}

// Put request
func (ctl *UserController) Put() {
	id := ctl.Ctx.Input.Param(":id")
	ctl.URL = "/user/"
	if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
		if user, err := md.GetUserByID(idInt64); err == nil {
			if err := ctl.ParseForm(&user); err == nil {
				fmt.Println(user.DepartmentID)
				fmt.Println(user.PositionID)
				fmt.Println(user.GroupIDs)
				fmt.Println(user.TeamIDs)
				var upateField []string
				if departmentID, err := ctl.GetInt64("department"); err == nil {
					if department, err := md.GetDepartmentByID(departmentID); err == nil {
						user.Department = department
						upateField = append(upateField, "Department")
					}
				}
				groupIDsStr := ctl.GetStrings("group")
				var groupIDs []int64
				for _, el := range groupIDsStr {
					if idInt64, err := strconv.ParseInt(el, 10, 64); err == nil {
						groupIDs = append(groupIDs, idInt64)
					}
				}
				if len(groupIDs) > 0 {
					var groups []*md.Group
					for _, groupID := range groupIDs {
						if group, err := md.GetGroupByID(groupID); err == nil {
							groups = append(groups, group)
						}
					}
					user.Groups = groups
					upateField = append(upateField, "Groups")
				}
				if positionID, err := ctl.GetInt64("position"); err == nil {
					if position, err := md.GetPositionByID(positionID); err == nil {
						user.Position = position
						upateField = append(upateField, "Position")
					}
				}
				if err := md.UpdateUserByID(user); err == nil {
					ctl.Redirect(ctl.URL+id+"?action=detail", 302)
				}
			}
		}
	}
	ctl.Redirect(ctl.URL+id+"?action=edit", 302)
}

// Get request
func (ctl *UserController) Get() {
	ctl.PageName = "用户管理"
	ctl.URL = "/user/"
	ctl.Data["URL"] = ctl.URL

	action := ctl.Input().Get("action")
	switch action {
	case "create":
		ctl.Create()
	case "edit":
		ctl.Edit()
	case "detail":
		ctl.Detail()
	case "changepasswd":
		ctl.ChangePwd()
	default:
		ctl.GetList()
	}
	ctl.Data["PageName"] = ctl.PageName + "\\" + ctl.PageAction

}

// Post request
func (ctl *UserController) Post() {
	action := ctl.Input().Get("action")
	ctl.URL = "/user/"
	switch action {
	case "validator":
		ctl.Validator()
	case "table": //bootstrap table的post请求
		ctl.PostList()
	case "create":
		ctl.PostCreate()
	default:
		ctl.PostList()
	}
}

// Create get user create page
func (ctl *UserController) Create() {
	ctl.Data["Action"] = "create"
	ctl.Data["Readonly"] = false
	ctl.PageAction = "创建"

	ctl.Layout = "base/base.html"
	ctl.TplName = "user/user_form.html"
}

// Detail display user info
func (ctl *UserController) Detail() {
	//获取信息一样，直接调用Edit
	ctl.Edit()
	ctl.Data["Readonly"] = true
	ctl.Data["MenuSelfInfoActive"] = "active"
	ctl.Data["Action"] = "detail"
}

// GetList display user with list
func (ctl *UserController) GetList() {
	viewType := ctl.Input().Get("view")
	if viewType == "" || viewType == "table" {
		ctl.Data["ViewType"] = "table"
	}
	ctl.PageAction = "列表"
	ctl.Data["tableId"] = "table-user"
	ctl.Data["MenuUserActive"] = "active"
	ctl.Layout = "base/base_list_view.html"
	ctl.TplName = "user/user_list_search.html"
}

// Validator js valid
func (ctl *UserController) Validator() {
	recordID, _ := ctl.GetInt64("recordId")
	name := strings.TrimSpace(ctl.GetString("Name"))
	result := make(map[string]bool)
	obj, err := md.GetUserByName(name)
	if err != nil {
		result["valid"] = true
	} else {
		if obj.Name == name {
			if recordID == obj.ID {
				result["valid"] = true
			} else {
				result["valid"] = false
			}

		} else {
			result["valid"] = true
		}

	}
	ctl.Data["json"] = result
	ctl.ServeJSON()
}

// PostList post request json response
func (ctl *UserController) PostList() {
	query := make(map[string]string)
	fields := make([]string, 0, 0)
	sortby := make([]string, 0, 0)
	order := make([]string, 0, 0)
	offset, _ := ctl.GetInt64("offset")
	limit, _ := ctl.GetInt64("limit")
	name := strings.TrimSpace(ctl.Input().Get("Name"))
	if name != "" {
		query["Name"] = name
	}
	if result, err := ctl.userList(query, fields, sortby, order, offset, limit); err == nil {
		ctl.Data["json"] = result
	}
	ctl.ServeJSON()

}
func (ctl *UserController) userList(query map[string]string, fields []string, sortby []string, order []string, offset int64, limit int64) (map[string]interface{}, error) {

	var users []md.User
	paginator, users, err := md.GetAllUser(query, fields, sortby, order, offset, limit)
	result := make(map[string]interface{})
	if err == nil {

		// result["recordsFiltered"] = paginator.TotalCount
		tableLines := make([]interface{}, 0, 4)
		for _, user := range users {

			oneLine := make(map[string]interface{})
			oneLine["Name"] = user.Name
			oneLine["NameZh"] = user.NameZh
			if user.Department != nil {
				oneLine["Department"] = user.Department.Name
			} else {
				oneLine["Department"] = "-"
			}
			if user.Position != nil {
				oneLine["Position"] = user.Position.Name
			} else {
				oneLine["Position"] = "-"
			}

			oneLine["Email"] = user.Email
			oneLine["Mobile"] = user.Mobile
			oneLine["Tel"] = user.Tel
			oneLine["IsAdmin"] = user.IsAdmin
			oneLine["Active"] = user.Active
			oneLine["Qq"] = user.Qq
			oneLine["ID"] = user.ID
			oneLine["id"] = user.ID
			oneLine["Wechat"] = user.WeChat

			tableLines = append(tableLines, oneLine)
		}
		result["data"] = tableLines

		if jsonResult, er := json.Marshal(&paginator); er == nil {
			result["paginator"] = string(jsonResult)
			result["total"] = paginator.TotalCount
		}
	}
	return result, err
}

// ChangePwd change password
func (ctl *UserController) ChangePwd() {
	ctl.Data["MenuChangePwdActive"] = "active"
	ctl.Layout = "base/base.html"
	ctl.TplName = "user/user_change_password_form.html"
}

//PostCreate create user with post params
func (ctl *UserController) PostCreate() {

	user := new(md.User)
	if err := ctl.ParseForm(user); err == nil {
		fmt.Println(user.DepartmentID)
		fmt.Println(user.PositionID)
		fmt.Println(user.GroupIDs)
		fmt.Println(user.TeamIDs)
		fmt.Println(ctl.GetStrings("Group"))
		if deparentID, err := ctl.GetInt64("Department"); err == nil {
			if department, err := md.GetDepartmentByID(deparentID); err == nil {
				user.Department = department
			}
		}
		if positionID, err := ctl.GetInt64("Position"); err == nil {
			if position, err := md.GetPositionByID(positionID); err == nil {
				user.Position = position
			}
		}
		if id, err := md.AddUser(user); err == nil {
			ctl.Redirect(ctl.URL+strconv.FormatInt(id, 10)+"?action=detail", 302)

		}
	}

}

// Edit edit user info
func (ctl *UserController) Edit() {
	id := ctl.Ctx.Input.Param(":id")
	userInfo := make(map[string]interface{})
	if id != "" {
		if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
			if user, err := md.GetUserByID(idInt64); err == nil {
				ctl.PageAction = user.Name + "(" + user.NameZh + ")"
				userInfo["ID"] = user.ID
				userInfo["Name"] = user.Name
				userInfo["NameZh"] = user.NameZh
				userInfo["Email"] = user.Email
				userInfo["Mobile"] = user.Mobile
				userInfo["Qq"] = user.Qq
				userInfo["Wechat"] = user.WeChat
				userInfo["Tel"] = user.Tel
				department := make(map[string]string)
				if user.Department != nil {
					department["ID"] = strconv.FormatInt(user.Department.ID, 10)
					department["Name"] = user.Department.Name
					userInfo["Department"] = department
				}
				groups := make([]interface{}, 0, 4)
				for _, group := range user.Groups {
					oneLine := make(map[string]interface{})
					oneLine["ID"] = group.ID
					oneLine["Name"] = group.Name
					oneLine["Description"] = group.Description
					groups = append(groups, oneLine)
				}
				userInfo["Groups"] = groups
				position := make(map[string]string)
				if user.Position != nil {
					position["ID"] = strconv.FormatInt(user.Position.ID, 10)
					position["Name"] = user.Position.Name
					userInfo["Position"] = position
				}
			}
		}
	}
	ctl.Data["RecordID"] = id
	ctl.Data["Action"] = "edit"
	ctl.Data["User"] = userInfo
	ctl.Layout = "base/base.html"
	ctl.TplName = "user/user_form.html"
}