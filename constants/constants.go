package constants

var (
	//超级管理员角色ID
	SuperAdminRoleID = ""
	//平台租户ID
	PlatformTenantID = ""
)

func SetSuperAdminRoleID(roleID string) {
	SuperAdminRoleID = roleID
}

func SetPlatformTenantID(tenantID string) {
	PlatformTenantID = tenantID
}
