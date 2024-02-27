package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

var server = "118.31.103.239"
var port = 1433
var user = "sa"
var password = "aabbcc1228.."
var database = "ZLBase"
var table = "sys_activity_switch"
var tableDesc = "活动开关配置"

type EntityInfo struct {
	TableName string
	TableDesc string
	FieldMap  map[string]string
}

func main() {
	entityInfo := &EntityInfo{
		TableName: table,
		TableDesc: tableDesc,
	}
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", server, user, password, port, database)
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	query := "SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = @tableName"
	rows, err := db.QueryContext(ctx, query, sql.Named("tableName", table))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fieldMap := make(map[string]string)
	for rows.Next() {
		var columnName, dataType string
		err := rows.Scan(&columnName, &dataType)
		if err != nil {
			log.Fatal(err)
		}

		// Map SQL data types to Golang data types
		var goType string
		switch dataType {
		case "int", "smallint", "tinyint":
			goType = "int32"
		case "nvarchar", "varchar", "nchar":
			goType = "string"
		case "datetime":
			goType = "time.Time"
		// Add more cases for other data types as needed
		default:
			goType = "interface{}"
		}

		fieldMap[columnName] = goType
	}
	entityInfo.FieldMap = fieldMap
	entityInfo.GenerateModel()
	entityInfo.GenerateLogic()
	entityInfo.GenerateParam()
	entityInfo.GenerateHandler()
}

// getFieldName 大驼峰命令
func getFieldName1(name string) string {
	s2 := strings.Split(name, "_")
	for k, v := range s2 {
		s2[k] = strings.ToUpper(string(v[0])) + v[1:]
	}
	return strings.Join(s2, "")
}

// getFieldName 小驼峰命令
func getFieldName2(name string) string {
	s2 := strings.Split(name, "_")
	for k, v := range s2 {
		if k != 0 {
			s2[k] = strings.ToUpper(string(v[0])) + v[1:]
		}

	}
	return strings.Join(s2, "")
}

func (e *EntityInfo) FilterField(field string) bool {
	if field == "created_at" || field == "updated_at" || field == "id" {
		return true
	}
	return false
}

func (e *EntityInfo) GenerateHandler() {
	entityName := getFieldName1(e.TableName)
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("package %s", e.TableName))
	buffer.WriteString("\n\n")

	buffer.WriteString(`import (
		"BlindBoxAdmin/cmd/box_admin/response"
		"BlindBoxAdmin/cmd/box_admin/util"
		"strconv"
	
		"github.com/gin-gonic/gin"
	)
	`)

	buffer.WriteString("\n\n")

	//查询方法
	buffer.WriteString(fmt.Sprintf(`
// BindApiRouter 绑定API路由
func BindApiRouter(api *gin.RouterGroup) {
	api.GET("system/%s", Get%ss)
	api.POST("system/%s", Post%s)
	api.PATCH("system/%s/:id", Patch%s)
	api.DELETE("system/%s/:id", Delete%s)
}
	`, e.TableName, entityName, e.TableName, entityName, e.TableName, entityName, e.TableName, entityName))

	//列表方法
	buffer.WriteString(fmt.Sprintf(`
	// Get%ss
	//	@Summary		获取%s列表
	//	@Description	获取%s列表
	//	@Tags			系统管理-%s模块
	//	@security		ApiKeyAuth
	//	@Param			Authorization	header	string								true	"设置到HTTP头内"
	//	@Param			data			query	Get%sReq							true	"编码JSON格式"
	//	@Success		200				object	response.Success{data=Get%sResp}	"A successful response."
	//	@Response		default			object	response.Failed						"An unexpected error response."
	//	@Router			/system/%s [get]
	func Get%ss(c *gin.Context) {
		var params = Get%sReq{}
		if err := util.BeforeApiLogic(c, &params); err != nil {
			_ = c.Error(err)
			return
		}
		util.PaginateParams(&params)
		resp, err := logicGet%ss(&params)
		if err != nil {
			response.Fail(c, err)
			_ = c.Error(err)
			return
		}
		response.Succ(c, resp)
	}`, entityName, e.TableDesc, e.TableDesc, e.TableDesc, entityName, entityName, e.TableName, entityName, entityName, entityName,
	))

	//新增方法
	buffer.WriteString(fmt.Sprintf(`
	// Post%s
	//	@Summary		添加%s
	//	@Description	添加%s
	//	@Tags			系统管理-%s模块
	//	@security		ApiKeyAuth
	//	@Param			Authorization	header	string								true	"设置到HTTP头内"
	//	@Param			data			body	Post%sReq							true	"编码JSON格式"
	//	@Success		200				object	response.Success{data=Post%sResp}	"A successful response."
	//	@Response		default			object	response.Failed						"An unexpected error response."
	//	@Router			/system/%s [post]
	func Post%s(c *gin.Context) {
		var params = Post%sReq{}
		if err := util.BeforeApiLogic(c, &params); err != nil {
			_ = c.Error(err)
			return
		}
	
		reqUserId := util.GetReqUserId(c)
		resp, err := logicAdd%s(reqUserId, &params)
		if err != nil {
			response.Fail(c, err)
			_ = c.Error(err)
			return
		}
		response.Succ(c, resp)
	}
`, entityName, tableDesc, tableDesc, tableDesc, entityName, entityName, e.TableName, entityName, entityName, entityName))

	//更新方法
	buffer.WriteString(fmt.Sprintf(`
// Patch%s
//	@Summary		更新%s
//	@Description	更新指定的%s
//	@Tags			系统管理-%s模块
//	@security		ApiKeyAuth
//	@Param			Authorization	header	string									true	"设置到HTTP头内"
//	@Param			id				path	integer									true	"角色ID"
//	@Param			data			body	Patch%sReq							true	"编码JSON格式"
//	@Success		200				object	response.Success{data=Patch%sResp}	"A successful response."
//	@Response		default			object	response.Failed							"An unexpected error response."
//	@Router			/system/%s/{id} [patch]
func Patch%s(c *gin.Context) {
	var params = Patch%sReq{}
	if err := util.BeforeApiLogic(c, &params); err != nil {
		_ = c.Error(err)
		return
	}

	reqUserId := util.GetReqUserId(c)
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := logicPatch%s(reqUserId, int32(id), &params)
	if err != nil {
		response.Fail(c, err)
		_ = c.Error(err)
		return
	}
	response.Succ(c, resp)
}`, entityName, e.TableDesc, e.TableDesc, e.TableDesc, entityName, entityName, e.TableName, entityName, entityName, entityName))

	// 更新方法
	buffer.WriteString(fmt.Sprintf(`
// Delete%s
//	@Summary		删除%s
//	@Description	删除指定的%s 
//	@Tags			系统管理-%s模块
//	@security		ApiKeyAuth
//	@Param			Authorization	header	string									true	"设置到HTTP头内"
//	@Param			id				path	integer									true	"角色ID"
//	@Success		200				object	response.Success{data=Delete%sResp}	"A successful response."
//	@Response		default			object	response.Failed							"An unexpected error response."
//	@Router			/system/%s/{id} [delete]
func Delete%s(c *gin.Context) {
	reqUserId := util.GetReqUserId(c)
	id, _ := strconv.Atoi(c.Param("id"))
	resp, err := logicDel%s(reqUserId, int32(id))
	if err != nil {
		response.Fail(c, err)
		_ = c.Error(err)
		return
	}
	response.Succ(c, resp)
}`, entityName, e.TableDesc, e.TableDesc, e.TableDesc, entityName, e.TableName, entityName, entityName))

	fmt.Println(buffer.String())

	os.WriteFile("handler.go", buffer.Bytes(), 0644)
}

func (e *EntityInfo) GenerateParam() {

	entityName := getFieldName1(e.TableName)

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("package %s", e.TableName))
	buffer.WriteString("\n\n\n")

	buffer.WriteString(fmt.Sprintf("// Get%sReq 请求-获取%s列表\n", entityName, e.TableDesc))
	buffer.WriteString(fmt.Sprintf("type Get%sReq struct {\n", entityName))
	buffer.WriteString("	Page      int    `form:\"page\"`      // 查询起始页\n")
	buffer.WriteString("	Limit     int    `form:\"limit\"`     // 每页记录数\n")
	buffer.WriteString("//TODO 查询字段\n")
	buffer.WriteString("}\n\n")

	buffer.WriteString(fmt.Sprintf("// Get%sResp 响应-获取%s列表\n", entityName, e.TableDesc))
	buffer.WriteString(fmt.Sprintf("type Get%sResp struct {\n", entityName))
	buffer.WriteString("Total int64 `json:\"total\"` // 记录总数\n")
	buffer.WriteString(fmt.Sprintf("Items []*db_%s.%s `json:\"items\"` // %s列表\n", e.TableName, entityName, e.TableDesc))
	buffer.WriteString("}\n\n")

	buffer.WriteString(fmt.Sprintf("// Post%sReq 请求-添加%s\n", entityName, e.TableDesc))
	buffer.WriteString(fmt.Sprintf("type Post%sReq struct{\n", entityName))
	for k, v := range e.FieldMap {
		if e.FilterField(k) {
			continue
		}
		buffer.WriteString(fmt.Sprintf("\t%s %s  `json:\"%s\" binding:\"required\"`\n", getFieldName1(k), v, getFieldName2(k)))
	}
	buffer.WriteString("}\n\n")

	buffer.WriteString(fmt.Sprintf("// Post%sResp 响应-添加%s\n", entityName, e.TableDesc))
	buffer.WriteString(fmt.Sprintf("type Post%sResp struct{\n", entityName))
	buffer.WriteString("}\n\n")

	buffer.WriteString(fmt.Sprintf("// Patch%sReq 请求-更新%s\n", entityName, e.TableDesc))
	buffer.WriteString(fmt.Sprintf("type Patch%sReq struct{\n", entityName))
	for k, v := range e.FieldMap {
		if e.FilterField(k) {
			continue
		}
		buffer.WriteString(fmt.Sprintf("\t%s *%s  `json:\"%s\" binding:\"required\"`\n", getFieldName1(k), v, getFieldName2(k)))
	}
	buffer.WriteString("}\n\n")

	buffer.WriteString(fmt.Sprintf("// Patch%sResp 响应-更新%s\n", entityName, e.TableDesc))
	buffer.WriteString(fmt.Sprintf("type Patch%sResp struct{\n", entityName))
	buffer.WriteString("}\n\n")

	buffer.WriteString(fmt.Sprintf("// Delete%sResp 响应-删除%s\n", entityName, e.TableDesc))
	buffer.WriteString(fmt.Sprintf("type Delete%sResp struct{\n", entityName))
	buffer.WriteString("}\n\n")

	fmt.Println(buffer.String())

	os.WriteFile("param.go", buffer.Bytes(), 0644)
}

func (e *EntityInfo) GenerateLogic() {

	entityName := getFieldName1(e.TableName)

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("package %s", e.TableName))
	buffer.WriteString("\n\n")

	buffer.WriteString(`import (
	"BlindBoxAdmin/cmd/box_admin/util"	
	"BlindBoxAdmin/internal/db/models/admin/log"	
	"BlindBoxAdmin/internal/db/util"
	"BlindBoxAdmin/internal/log"
	)
	`)

	buffer.WriteString("\n\n")

	//查询方法
	buffer.WriteString(fmt.Sprintf(`
	// logicGet%ss 获取%s列表的逻辑
	func logicGet%ss(req *Get%sReq) (*Get%sResp, error) {
		var whereStr string
		var whereValues = make([]interface{}, 0)
		//TODO 查询
		// if req.UserId > 0 { // first param
		// 	whereStr += "user_id=?"
		// 	whereValues = append(whereValues, req.UserId)
		// }
	
		whereScopes := db_util.Where(whereStr, whereValues...)

		items, total, err := db_%s.Query%sList(db_util.Paginate(req.Page, req.Limit),whereScopes)
		if err != nil {
			log.Error("Query%sList failed", "error", err)
			return nil, util.ErrorWithRequestFailed("")
		}
		return &Get%sResp{Total: total, Items: items}, nil
	}
	`, entityName, tableDesc, entityName, entityName, entityName, e.TableName, entityName, entityName, entityName))

	addEntityItemsStr := ""
	for k, _ := range e.FieldMap {
		if e.FilterField(k) {
			continue
		}
		addEntityItemsStr += fmt.Sprintf("%s:req.%s, \n", getFieldName1(k), getFieldName1(k))
	}

	//添加方法
	buffer.WriteString(fmt.Sprintf(`
// logicAdd%s 添加%s的逻辑
func logicAdd%s(reqUserId int32, req *Post%sReq) (*Post%sResp, error) {
	//TODO 判断 是否存在 
	// isExist, err := db_%s.CheckSysNameExist(req.Name)
	// if err != nil {
	// 	log.Error("Check%sExist failed", "error", err)
	// 	return nil, util.ErrorWithRequestFailed("")
	// }
	// if isExist {
	// 	return nil, util.ErrorWithRequestFailed("%s已存在")
	// }

	item := &db_%s.%s{
		%s
	}
	if err := db_%s.Add%sData(item); err != nil {
		log.Error("Add%sData failed", "error", err)
		return nil, util.ErrorWithRequestFailed("")
	}

	
	if err := db_sys_log.AddSystemOperateLog(4, reqUserId, int64(item.Id), "添加%s");err != nil {
		log.Error("AddSystemOperateLog failed", "error", err)
	}

	return &Post%sResp{}, nil
}
	`, entityName, e.TableDesc, entityName, entityName, entityName, e.TableName, e.TableName, e.TableName,
		e.TableName, entityName, addEntityItemsStr, e.TableName, entityName, entityName, e.TableDesc, entityName,
	))

	updateEntityItemsStr := ""
	for k, v := range e.FieldMap {
		if e.FilterField(k) {
			continue
		}
		if v == "string" {
			updateEntityItemsStr += fmt.Sprintf(`util.UpdateFieldWithString(updateItems, "%s", req.%s, item.%s)`, k, getFieldName1(k), getFieldName1(k))
			updateEntityItemsStr += "\n"
		} else {
			updateEntityItemsStr += fmt.Sprintf(`util.UpdateFieldWithInt32(updateItems, "%s", req.%s, item.%s)`, k, getFieldName1(k), getFieldName1(k))
			updateEntityItemsStr += "\n"
		}

	}

	//更新方法
	buffer.WriteString(fmt.Sprintf(`
// logicPatch%s 更新%s的逻辑
func logicPatch%s(reqUserId int32, itemId int32, req *Patch%sReq) (*Patch%sResp, error) {
	item, err := db_%s.Query%sData(itemId)
	if err != nil {
		log.Error("Query%sData failed", "error", err)
		return nil, util.ErrorWithRequestFailed("")
	}
	if item == nil {
		return nil, util.ErrorWithRequestFailed("更新的%s不存在")
	}

	updateItems := make(db_util.Values)
 

	%s

	if len(updateItems) > 0 {
		if err = db_%s.Update%sData(itemId, updateItems); err != nil {
			log.Error("Update%sData failed", "error", err)
			return nil, util.ErrorWithRequestFailed("")
		}

		err = db_sys_log.AddSystemOperateLog(4, reqUserId, int64(itemId), "更新%s")
		if err != nil {
			log.Error("AddSystemOperateLog failed", "error", err)
		}
	}

	return &Patch%sResp{}, nil
}
`, entityName, tableDesc, entityName, entityName, entityName, e.TableName, entityName, entityName, e.TableDesc, updateEntityItemsStr, e.TableName, entityName, entityName, e.TableDesc, entityName))

	//删除方法
	buffer.WriteString(fmt.Sprintf(`
	// logicDel%s 删除%s的逻辑
	func logicDel%s(reqUserId int32, itemId int32) (*Delete%sResp, error) { 
		if err := db_%s.Delete%sData(itemId); err != nil {
			log.Error("Delete%sData failed", "error", err)
			return nil, util.ErrorWithRequestFailed("")
		}
	
		err := db_sys_log.AddSystemOperateLog(4, reqUserId, int64(itemId), "删除%s")
		if err != nil {
			log.Error("AddSystemOperateLog failed", "error", err)
		}
	
		return &Delete%sResp{}, nil
	}	`, entityName, e.TableDesc, entityName, entityName, e.TableName, entityName, entityName, e.TableDesc, entityName))

	fmt.Println(buffer.String())

	os.WriteFile("logic.go", buffer.Bytes(), 0644)
}

func (e *EntityInfo) GenerateModel() {

	entityName := getFieldName1(e.TableName)

	var buffer bytes.Buffer

	// Generate Golang struct based on table columns
	// fmt.Printf("type %s struct {\n", getFieldName(table))
	buffer.WriteString(fmt.Sprintf("package db_%s", e.TableName))
	buffer.WriteString("\n\n")

	buffer.WriteString(`import (
		"BlindBoxAdmin/internal/db/mssql"
		db_util "BlindBoxAdmin/internal/db/util"
		"errors"
	
		"gorm.io/gorm"
	)
	`)

	buffer.WriteString("\n\n")
	buffer.WriteString(fmt.Sprintf("type %s struct {\n", entityName))

	for k, v := range e.FieldMap {
		buffer.WriteString(fmt.Sprintf("\t%s %s  `gorm:\"column:%s\" json:\"%s\"`\n", getFieldName1(k), v, k, getFieldName2(k)))
	}

	buffer.WriteString("}\n\n")

	//查询方法
	buffer.WriteString(fmt.Sprintf(`
// Query%sList 查询%s列表
func Query%sList(paginate,where db_util.ScopesFunc) ([]*%s, int64, error) {
	var mydb = mssql.GetDB()
	tx := mydb.Table("%s")
	if where != nil {
		tx = tx.Scopes(where)
	}

	var count = int64(0)
	if err := tx.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var items []*%s	
	if err := tx.Scopes(paginate).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, count, nil
}
	`, entityName, tableDesc, entityName, entityName, e.TableName, entityName))

	//添加方法
	buffer.WriteString(fmt.Sprintf(`
// Add%sData 添加新的%s
func Add%sData(new *%s) error {
	var mydb = mssql.GetDB()
	new.CreatedAt = time.Now()
	new.UpdatedAt = time.Now()
	tx := mydb.Table("%s")
	if err := tx.Create(new).Error; err != nil {
		return err
	}
	return nil
}
	`, entityName, tableDesc, entityName, entityName, e.TableName))

	//查询方法
	buffer.WriteString(fmt.Sprintf(`
// Query%sData 查询指定的%s数据
func Query%sData(id int32) (*%s, error) {
	var mydb = mssql.GetDB()
	var item = %s{}
	tx := mydb.Table("%s").Where("id=?", id)
	if err := tx.Take(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}
			`, entityName, tableDesc, entityName, entityName, entityName, e.TableName))

	//更新方法
	buffer.WriteString(fmt.Sprintf(`
// Update%sData 更新指定的%s数据
func Update%sData(id int32, values map[string]interface{}) error {
	var mydb = mssql.GetDB()
	values["updated_at"] = time.Now()
	tx := mydb.Table("%s").Where("id=?", id)
	if err := tx.Updates(values).Error; err != nil {
		return err
	}
	return nil
}
				`, entityName, tableDesc, entityName, e.TableName))

	//删除方法
	buffer.WriteString(fmt.Sprintf(`
// Delete%sData 删除%s数据
func Delete%sData(id int32) error {
	var mydb = mssql.GetDB()
	tx := mydb.Table("%s").Where("id=?", id)
	if err := tx.Delete(nil).Error; err != nil {
		return err
	}
	return nil
}
					`, entityName, tableDesc, entityName, e.TableName))

	fmt.Println(buffer.String())

	os.WriteFile(e.TableName+".go", buffer.Bytes(), 0644)
}
