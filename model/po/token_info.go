package po

type TokenInfo struct {
	Base
	Token string `gorm:"type:varchar(100)"`
	Ip    string `gorm:"type:varchar(256)"`
	Port  string `gorm:"type:char(5)"`
}
