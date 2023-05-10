package keys

const TargetGroupKey = "TARGET_GROUPV2#"

type targetGroupKeys struct {
	PK1 string
	SK1 func(id string) string
}

var TargetGroup = targetGroupKeys{
	PK1: TargetGroupKey,
	SK1: func(id string) string { return id + "#" },
}
