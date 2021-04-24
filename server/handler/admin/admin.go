package admin

import (
	"github.com/gin-gonic/gin"
	"labsystem/configs"
	"labsystem/model"
	"labsystem/model/srverr"
	"labsystem/server/handler"
	adminSrv "labsystem/service/admin"
	"labsystem/util/rsa"
	"net/http"
)

type HandlerAdmin struct {
	Srv     adminSrv.ServiceAdmin
	handles []*handler.Handle
}

var profilePath = configs.CurProjectPath() + "/static/profile/"

const fileMax = 1024 * 1024 * 4

func (h *HandlerAdmin) verifyAdmin(ctx *gin.Context) *model.Admin {
	uid, ok := ctx.Keys["uid"]
	if !ok {
		return nil
	}
	if rid, ok := ctx.Keys["rid"]; !ok || int(rid.(float64)) != model.Administrator.Int() {
		return nil
	}
	admin := h.Srv.QueryAdminById(uint(uid.(float64)))
	if admin == nil {
		return nil
	}

	return admin
}

func (h *HandlerAdmin) RegisterAdminHandles(rg *gin.RouterGroup, authRg *gin.RouterGroup) {
	// the RouterGroup rg mustn't be Authenticate
	{
		rg.POST("/login", h.login)
	}
	// the RouterGroup hRg must be Authenticate
	{
		authRg.POST("/info", h.adminInfo)
		authRg.POST("/list", h.adminList)
		authRg.POST("/create", h.createAdmin)
		authRg.POST("/update", h.updateAdmin)
		authRg.POST("/delete", h.deleteAdmin)
		authRg.POST("/class/create", h.createClass)
		authRg.POST("/teacher/create", h.createTeacher)
		authRg.POST("/user/list", h.userList)
		authRg.POST("/class/list", h.classList)
		authRg.POST("/user/delete", h.deleteUsers)
	}
}

func (h *HandlerAdmin) login(ctx *gin.Context) {
	var req loginReq
	if err := ctx.Bind(&req); err != nil || !req.Valid() {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrInvalidParams, nil))
		return
	}
	t, err := h.Srv.Login(req.AdminNick, req.Password, req.Key, req.VCode)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(err, nil))
		return
	}

	ctx.JSON(http.StatusOK, handler.NewResp(nil, t))
}

func (h *HandlerAdmin) adminInfo(ctx *gin.Context) {
	admin := h.verifyAdmin(ctx)
	if admin == nil {
		ctx.JSON(http.StatusForbidden, handler.NewResp(srverr.ErrForbidden, nil))
		return
	}
	resp := new(InfoResp)
	resp.Name = admin.NickName
	resp.Powers = make([]*PowerOwner, len(model.PowerList))
	resp.Id = admin.ID
	for i, v := range model.PowerList {
		resp.Powers[i] = &PowerOwner{
			Name:  v.Name,
			Power: v.No,
			Own:   admin.Power.Own(v.No),
		}
	}

	ctx.JSON(http.StatusOK, handler.NewResp(nil, resp))
}

func (h *HandlerAdmin) adminList(ctx *gin.Context) {
	admin := h.verifyAdmin(ctx)
	if admin == nil {
		ctx.JSON(http.StatusForbidden, handler.NewResp(srverr.ErrForbidden, nil))
		return
	}
	var req *ListReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrInvalidParams, nil))
		return
	}
	list, totalPage, totalCount := h.Srv.AdminList(&adminSrv.ListOpt{
		CreatedBy: req.CreatedBy,
	}, req.Page, req.PageSize)
	items := make([]*Item, len(list))
	for i, v := range list {
		items[i] = new(Item)
		items[i].ID = v.ID
		items[i].Name = v.NickName
		if creator := h.Srv.QueryAdminById(v.CreatedBy); creator != nil {
			items[i].CreatedBy = creator.NickName

		}
		items[i].CreatedAt = v.CreatedAt
		items[i].Power = make([]*PowerOwner, len(model.PowerList))
		for k, p := range model.PowerList {
			items[i].Power[k] = new(PowerOwner)
			items[i].Power[k].Power = p.No
			items[i].Power[k].Name = p.Name
			if v.Power.Own(p.No) {
				items[i].Power[k].Own = true
			}
		}
	}
	ctx.JSON(http.StatusOK, handler.NewResp(nil, &ListResp{
		List:       items,
		TotalPage:  totalPage,
		TotalCount: totalCount,
	}))
}

func (h *HandlerAdmin) createAdmin(ctx *gin.Context) {
	admin := h.verifyAdmin(ctx)
	if admin == nil {
		ctx.JSON(http.StatusForbidden, handler.NewResp(srverr.ErrForbidden, nil))
		return
	}
	// verify power
	if !admin.Power.Own(model.PowerCreateAdmin) {
		ctx.JSON(http.StatusUnauthorized, handler.NewResp(srverr.ErrOwnPower, nil))
		return
	}
	var req *CreateAdminReq
	if err := ctx.BindJSON(&req); err != nil || !req.Valid() {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrInvalidParams, nil))
		return
	}
	// verify power legitimacy
	p, _ := model.IntToPower(req.Power)
	if !admin.Power.Own(p) {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrInvalidParams, nil))
		return
	}
	if err := h.Srv.CreateAdmin(&model.Admin{
		NickName:  req.Name,
		Password:  req.Password,
		Power:     p,
		CreatedBy: admin.ID,
	}); err != nil {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrSystemException, nil))
		return
	}

	ctx.JSON(http.StatusOK, handler.NewResp(nil, nil))
}

func (h *HandlerAdmin) createTeacher(ctx *gin.Context) {
	admin := h.verifyAdmin(ctx)
	if admin == nil {
		ctx.JSON(http.StatusForbidden, handler.NewResp(srverr.ErrForbidden, nil))
		return
	}
	// verify power
	if !admin.Power.Own(model.PowerCreateTeacher) {
		ctx.JSON(http.StatusUnauthorized, handler.NewResp(srverr.ErrOwnPower, nil))
		return
	}
	req := &CreateTeacherReq{
		UserNo:   ctx.Request.FormValue("user_no"),
		RealName: ctx.Request.FormValue("real_name"),
		Password: ctx.Request.FormValue("password"),
		Class:    ctx.Request.FormValue("class"),
		FileName: ctx.Request.FormValue("file_name"),
	}
	if !req.Valid() {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrInvalidParams, nil))
		return
	}
	_, file, err := ctx.Request.FormFile(req.FileName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, handler.NewResp(srverr.ErrUpload, nil))
		return
	}
	if file.Size > fileMax {
		ctx.JSON(http.StatusInternalServerError, handler.NewResp(srverr.ErrFileMax, nil))
		return
	}
	err = ctx.SaveUploadedFile(file, profilePath+req.FileName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, handler.NewResp(srverr.ErrUpload, nil))
		return
	}

	err = h.Srv.CreateTeacher(&model.User{
		UserNo:     req.UserNo,
		RealName:   req.RealName,
		Password:   req.Password,
		Class:      req.Class,
		ProfileUrl: profilePath + req.FileName,
		CreatedBy:  admin.NickName,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(err, nil))
		return
	}

	ctx.JSON(http.StatusOK, handler.NewResp(nil, nil))
}

func (h *HandlerAdmin) updateAdmin(ctx *gin.Context) {
	admin := h.verifyAdmin(ctx)
	if admin == nil {
		ctx.JSON(http.StatusForbidden, handler.NewResp(srverr.ErrForbidden, nil))
		return
	}
	var req *UpdateAdminReq
	if err := ctx.BindJSON(&req); err != nil || !req.Valid() {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrInvalidParams, nil))
		return
	}
	if ok := rsa.Compare(admin.Password, req.OldPassword); !ok {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrInvalidParams, nil))
		return
	}
	ok := h.Srv.UpdateAdmin(admin.ID, req.Name, req.NewPassword)
	if !ok {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrUpdateFailed, nil))
		return
	}

	ctx.JSON(http.StatusOK, handler.NewResp(nil, nil))
}

func (h *HandlerAdmin) deleteAdmin(ctx *gin.Context) {
	admin := h.verifyAdmin(ctx)
	if admin == nil {
		ctx.JSON(http.StatusForbidden, handler.NewResp(srverr.ErrForbidden, nil))
		return
	}
	var req *DeleteAdminReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrInvalidParams, nil))
		return
	}
	if ok := h.Srv.DeleteAdmin(admin.ID, req.ID); !ok {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrDeleteFailed, nil))
		return
	}

	ctx.JSON(http.StatusOK, handler.NewResp(nil, nil))
}

func (h *HandlerAdmin) createClass(ctx *gin.Context) {
	admin := h.verifyAdmin(ctx)
	if admin == nil {
		ctx.JSON(http.StatusForbidden, handler.NewResp(srverr.ErrForbidden, nil))
		return
	}
	// verify power
	if !admin.Power.Own(model.PowerCreateClass) {
		ctx.JSON(http.StatusUnauthorized, handler.NewResp(srverr.ErrOwnPower, nil))
		return
	}
	var req *CreateClassReq
	if err := ctx.BindJSON(&req); err != nil || !req.Valid() {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrInvalidParams, nil))
		return
	}
	h.Srv.CreateClass(&model.Class{
		ClassNo: req.ClassNo,
	})

	ctx.JSON(http.StatusOK, handler.NewResp(nil, nil))
}

func (h *HandlerAdmin) userList(ctx *gin.Context) {
	admin := h.verifyAdmin(ctx)
	if admin == nil {
		ctx.JSON(http.StatusForbidden, handler.NewResp(srverr.ErrForbidden, nil))
		return
	}
	var req *UserListReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrInvalidParams, nil))
		return
	}
	list, totalPage, totalCount := h.Srv.UserList(req.Page, req.PageSize)
	if list == nil {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrSystemException, nil))
		return
	}
	resp := new(UserListResp)
	resp.List = make([]*UserItem, len(list))
	for i, v := range list {
		resp.List[i] = new(UserItem)
		resp.List[i].ID = v.ID
		resp.List[i].UserNo = v.UserNo
		resp.List[i].RealName = v.RealName
		resp.List[i].Class = v.Class
		resp.List[i].Status = v.Status
		resp.List[i].ProfileUrl = v.ProfileUrl
		resp.List[i].CreatedAt = v.CreatedAt
		resp.List[i].CreatedBy = v.CreatedBy
	}
	resp.TotalPage = totalPage
	resp.TotalCount = totalCount

	ctx.JSON(http.StatusOK, handler.NewResp(nil, resp))
}

func (h *HandlerAdmin) classList(ctx *gin.Context) {
	admin := h.verifyAdmin(ctx)
	if admin == nil {
		ctx.JSON(http.StatusForbidden, handler.NewResp(srverr.ErrForbidden, nil))
		return
	}
	var req *ClassListReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrInvalidParams, nil))
		return
	}
	list, totalPage, totalCount := h.Srv.ClassList(req.Page, req.PageSize)
	if list == nil {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrSystemException, nil))
		return
	}
	resp := new(ClassListResp)
	resp.List = make([]*ClassItem, len(list))
	for i, v := range list {
		resp.List[i] = new(ClassItem)
		resp.List[i].ID = v.ID
		resp.List[i].ClassNo = v.ClassNo
		resp.List[i].CreatedAt = v.CreatedAt
	}
	resp.TotalCount, resp.TotalPage = totalCount, totalPage
	ctx.JSON(http.StatusOK, handler.NewResp(nil, resp))
}

func (h *HandlerAdmin) deleteUsers(ctx *gin.Context) {
	admin := h.verifyAdmin(ctx)
	if admin == nil {
		ctx.JSON(http.StatusForbidden, handler.NewResp(srverr.ErrForbidden, nil))
		return
	}
	if !admin.Power.Own(model.PowerDeleteUser) {
		ctx.JSON(http.StatusUnauthorized, handler.NewResp(srverr.ErrOwnPower, nil))
	}
	var req *DeleteUsersReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrInvalidParams, nil))
		return
	}
	if !h.Srv.DeleteUsers(req.Ids) {
		ctx.JSON(http.StatusBadRequest, handler.NewResp(srverr.ErrDeleteFailed, nil))
	}

	ctx.JSON(http.StatusOK, handler.NewResp(nil, nil))
}