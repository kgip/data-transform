package constant

import "time"

const (
	ReturnRedPacketBalanceTopic = "redpacket"
	RedPacketExpireDuration     = time.Hour * 24 * 7 // 7天
	MinSingleRedPacketBalance   = 0.01               //红包最少金额
)
