package model

type PlatformTenant struct {
	CommonBase
	Id           uint64 `gorm:"type:int(11);primary_key;AUTO_INCREMENT;column:id"`
	PlatformCode string `gorm:"type:varchar(255);column:platform_code"`
	Prefix       string `gorm:"type:varchar(255);column:prefix"`
	Host         string `gorm:"type:varchar(255);column:host"`
	TenantCode   string `gorm:"type:varchar(255);column:tenant_code"`
}

func (PlatformTenant) TableName() string {
	return "g_platform_tenant"
}
