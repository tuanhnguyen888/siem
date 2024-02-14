package rules

import "strings"

func IsDiscoveryIam (event map[string]interface{})(b bool, alert string,level string) {
	if strings.Contains(event["event_name"].(string), "List") && strings.Contains(event["event_source"].(string),"iam.amazonaws.com") {
		return true, "Action List in service IAM" , a1
	}

	return false,"",""
}

func IsAssumeRole (event map[string]interface{})(b bool, alert string,level string) {
	if strings.Contains(event["event_name"].(string), "AssumeRole") && strings.Contains(event["event_source"].(string),"sts"){
		return true, "Execute AssumeRole in service STS", a1
	}
	return false,"",""
}

func IsBehavior (event map[string]interface{})(b bool, alert string,level string) {
	val, ok := event["event_name"].(string)
	if ok {
		switch val {
		case "GetPasswordData":
			return true,"Action Get Password Data of user", a1
		case "GetAccountPasswordPolicy":
			return true,"Action Get Account Password Policy in IAM", a1
		case "UpdateAccountPasswordPolicy":
			return true," Action Update Account Password Policy",a1

		}
	}
	return false,"",""
}

func IsActionKMS(event map[string]interface{})(b bool, alert string,level string)   {
	val, ok := event["event_name"].(string)
	if ok {
		switch val {
		case "DescribeKey":
			return true,"Action Describe Key KMS", a1
		case "CreateKey":
			return true,"Action Create Key KMS", a1
		}
	}
	return false,"",""
}

func IsActionEc2 (event map[string]interface{})(b bool, alert string,level string) {
	val, ok := event["event_name"].(string)
	if ok {
		switch val {
		case "StopInstance":
			return true,"Stop Instance EC2", a1
		case "CreateKey":
			return true,"Stop Instances EC2", a2
		case "TerminateInstances":
			return true,"Terminate Instances EC2", a2
		case "GetInstanceSnapshot":
			return true,"Get Instance Snapshot EC2", a1
		case "CopySnapshot":
			return true,"Copy Snapshot", a1
		case "CreateInstanceSnapshot":
			return true,"Create Instance Snapshot EC2", a1
		}
	}
	return false,"",""
}
