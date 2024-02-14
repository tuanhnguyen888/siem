package rules

import "strings"

func ISshFail(event map[string]interface{})  (b bool, alert string,level string)  {
	if ( event["spt"].(int) == int(22) || event["dst"].(int) == 22) && strings.Contains(event["event_name"].(string),"REJECT") {
		return true, "Connection Ssh Fail" , a1
	}
	return false,"",""
}
