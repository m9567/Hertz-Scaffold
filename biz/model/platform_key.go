package model

type PlatformKey struct {
	CommonBase
	Id uint64 `gorm:"type:int(11);primary_key;AUTO_INCREMENT;column:id"`
	//Agent    string `gorm:"type:varchar(255);column:agent;not null"`    //代理
	Code     string `gorm:"type:varchar(255);column:code;not null"`     // 场馆
	Currency string `gorm:"type:varchar(255);column:currency;not null"` //币种
	KeyJson  string `gorm:"type:text;column:key_json;not null"`         //key json

}

func (PlatformKey) TableName() string {
	return "g_platform_key"
}
