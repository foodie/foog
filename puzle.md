Id:   ((time.Now().UnixNano() / 1000000) << 22) | int64((appId&0x3ff)<<12) | (counter & 0xfff),
是什么情况？