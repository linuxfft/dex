package suc

type sucLoginAuthReq struct {
	Domain   string `json:"domain" validate:"required,lt=20"`
	Account  string `json:"account" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type sucLoginRole struct {
	RoleCode string `json:"roleCode"`
	RoleName string `json:"roleName"`
	RoleDesc string `json:"roleDesc"`
}

type sucLoginRespData struct {
	SecurityToken     string         `json:"securityToken"`
	PwdExpirationTime int            `json:"pwdExpirationTime"`
	Roles             []sucLoginRole `json:"roles"`
	Account           string         `json:"account"`
	Domain            string         `json:"domain"`
}

type sucLoginResp struct {
	Success bool             `json:"success"`
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    sucLoginRespData `json:"data"`
}
