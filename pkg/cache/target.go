package cache

type Target struct {
	ID          string
	AccessRules []string
	// @todo replace with detailed field
	Fields map[string]string
}

// func (d *Target) DDBKeys() (ddb.Keys, error) {
// 	keys := ddb.Keys{
// 		PK: keys.TargetGroupResource.PK1,
// 		SK: keys.TargetGroupResource.SK1(d.TargetGroupID, d.ResourceType, d.Resource.ID),
// 	}

// 	return keys, nil
// }
