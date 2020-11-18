package suc

type sucLoginAuthReq struct {
	Domain string `json:"domain" validate:"required,lt=20"`
	Account string `json:"account" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type sucLoginResp struct {
	Success           bool   `json:"success"`
	Code              int    `json:"code"`
	Message           string `json:"message"`
	Data              string `json:"data"`
	SecurityToken     string `json:"securityToken"`
	PwdExpirationTime string `json:"pwdExpirationTime"`
	Roles             string `json:"roles"`
	RoleCode          string `json:"roleCode"`
	RoleName          string `json:"roleName"`
	RoleDesc          string `json:"roleDesc"`
}
