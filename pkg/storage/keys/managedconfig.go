package keys

const ManagedProviderConfigKey = "MANAGED_DEPLOYMENT_CONFIG#PROVIDERS"
const ManagedNotificationsConfigKey = "MANAGED_DEPLOYMENT_CONFIG#NOTIFICATIONS"

type managedConfigKeys struct {
	PK string
	SK string
}

var ManagedProviderConfig = managedConfigKeys{
	PK: ManagedProviderConfigKey,
	SK: ManagedProviderConfigKey,
}

var ManagedNotificationConfig = managedConfigKeys{
	PK: ManagedNotificationsConfigKey,
	SK: ManagedNotificationsConfigKey,
}
