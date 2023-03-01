package keys

import "fmt"

const TargetRouteKey = "TARGET_ROUTE#"

type targetRouteKeys struct {
	PK1             string
	SK1             func(group string, handler string, kind string) string
	SK1GroupHandler func(group string, handler string) string
	SK1Group        func(group string) string
	GSI1PK          func(group string) string
	GSI1SK          func(valid bool, priority int) string
	GSI1SKValid     func(valid bool) string
	GSI2PK          func(handler string) string
	GSI2SK          func(group string) string
}

var TargetRoute = targetRouteKeys{
	PK1: TargetRouteKey,
	SK1: func(group string, handler string, kind string) string {
		return group + "#" + handler + "#" + kind + "#"
	},
	SK1GroupHandler: func(group string, handler string) string {
		return group + "#" + handler + "#"
	},
	SK1Group: func(group string) string { return group + "#" },
	GSI1PK: func(group string) string {
		return TargetRouteKey + "GROUP#" + group + "#"
	},
	GSI1SK: func(valid bool, priority int) string {
		return fmt.Sprintf("%v#%d#", valid, priority)
	},
	GSI1SKValid: func(valid bool) string {
		return fmt.Sprintf("%v#", valid)
	},
	GSI2PK: func(handler string) string {
		return TargetRouteKey + "HANDLER#" + handler + "#"
	},
	GSI2SK: func(group string) string {
		return group + "#"
	},
}
