package threatAnalysis

import (
	"correl/common"
	"correl/rules"
	"errors"
	"github.com/hashicorp/go-uuid"
	"github.com/sirupsen/logrus"
	"time"
)

func (f *FilterConfig) AnalysisEvent() error  {
	logrus.Info("Analysing Event")
	for {
		select {
		case event, ok := <-f.msgInput:
			if !ok {
				return errors.New("channel inputFilter closed. Exiting")
			}

			device, _ := event["device_product"].(string)
			switch device {
			case "AWS Cloudtrail":
				f.analysisAuditLog(event)
			case "AWS Route53":
				f.analysisVPC(event)
			case "AWS VPC Flow logs":
				f.analysisVPC(event)
			default:
				logrus.Errorf("Event type invalid")
			}

			logrus.Printf("DONE analysis event ")

		}
	}
}

func (f *FilterConfig) CreateAlert(evt map[string]interface{}, alertMsg string, level string)  {
	id, _ := uuid.GenerateUUID()
	alert := map[string]interface{}{
		"id": id,
		"message" : alertMsg,
		"event_list": []map[string]interface{}{evt},
		"level": level,
		"timestamp": time.Now().UnixMilli(),
	}
	logrus.Printf("Alert: %#v",alert)
	common.PushEvtToChan(alert, f.msgOutput)

}

func (f *FilterConfig) analysisAuditLog(event map[string]interface{})  {
	value, ok := event["src"].(string)
	if ok {
		b,alert,  level := rules.IsMaliciousIP(value)
		if b {
			f.CreateAlert(event,alert, level)
		}
	}

	if b,alert,  level := rules.IsDiscoveryIam(event) ; b {
		f.CreateAlert(event,alert, level)
	}

	if b,alert,  level := rules.IsActionKMS(event) ; b {
		f.CreateAlert(event,alert, level)
	}

	if b,alert,  level := rules.IsBehavior(event) ; b {
		f.CreateAlert(event,alert, level)
	}

	if b,alert,  level := rules.IsActionEc2(event) ; b {
		f.CreateAlert(event,alert, level)
	}

	//if b,alert,  level := rules.IsAssumeRole(event) ; b {
	//	f.CreateAlert(event,alert, level)
	//}

}

func (f *FilterConfig) analysisVPC(event map[string]interface{})  {
	value, ok := event["dst"].(string)
	if ok {
		b,alert,  level := rules.IsMaliciousIP(value)
		if b {
			f.CreateAlert(event,alert, level)
		}
	}

	value, ok = event["src"].(string)
	if ok {
		b,alert,  level := rules.IsMaliciousIP(value)
		if b {
			f.CreateAlert(event,alert, level)
		}
	}

	if b,alert,  level := rules.ISshFail(event) ; b {
		f.CreateAlert(event,alert, level)
	}
}

func (f *FilterConfig) analysisRoute53(event map[string]interface{})  {
	value, ok := event["query_name"].(string)
	if ok {
		b,alert,  level := rules.IsMaliciousURL(value)
		if b {
			f.CreateAlert(event,alert, level)
		}
	}
}