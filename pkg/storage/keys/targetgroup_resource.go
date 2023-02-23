package keys

const TargetGroupResourceKey = "TARGET_GROUP_RESOURCE#"

type targetGroupResourceKeys struct {
	PK1                    string
	SK1                    func(targetGroupID, resourceType string, value string) string
	SK1TargetGroupResource func(targetGroupID, resourceType string) string
}

var TargetGroupResource = targetGroupResourceKeys{
	PK1: TargetGroupResourceKey,
	SK1: func(targetGroupID, resourceType string, value string) string {
		return targetGroupID + "#" + resourceType + "#" + value
	},
	SK1TargetGroupResource: func(targetGroupID, resourceType string) string {
		return targetGroupID + "#" + resourceType + "#"
	},
}
