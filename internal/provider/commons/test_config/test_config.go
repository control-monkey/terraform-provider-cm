package test_config

import "os"

// TestConfig holds all test environment variables with CM_TEST_ prefix
type TestConfig struct {
	// Provider and repository configuration
	ProviderId string
	RepoName   string

	// Cloud and organization configuration
	CloudAccountId string
	OrgId          string

	// Control policy configuration
	ControlPolicyId      string
	ControlPolicyGroupId string

	// Notification configuration
	SlackWebhookUrl            string
	SlackAppId                 string
	NotificationEndpointEmail1 string
	NotificationEndpointEmail2 string
}

// Global test configuration instance
var Config *TestConfig

// init initializes the test configuration from environment variables
func init() {
	Config = &TestConfig{
		ProviderId:                 os.Getenv("CM_TEST_PROVIDER_ID"),
		RepoName:                   os.Getenv("CM_TEST_REPO_NAME"),
		CloudAccountId:             os.Getenv("CM_TEST_CLOUD_ACCOUNT_ID"),
		OrgId:                      os.Getenv("CM_TEST_ORG_ID"),
		ControlPolicyId:            os.Getenv("CM_TEST_CONTROL_POLICY_ID"),
		ControlPolicyGroupId:       os.Getenv("CM_TEST_CONTROL_POLICY_GROUP_ID"),
		SlackWebhookUrl:            os.Getenv("CM_TEST_SLACK_WEBHOOK_URL"),
		SlackAppId:                 os.Getenv("CM_TEST_SLACK_APP_ID"),
		NotificationEndpointEmail1: os.Getenv("CM_TEST_NOTIFICATION_ENDPOINT_EMAIL1"),
		NotificationEndpointEmail2: os.Getenv("CM_TEST_NOTIFICATION_ENDPOINT_EMAIL2"),
	}
}

// GetProviderId returns the test provider ID
func GetProviderId() string {
	return Config.ProviderId
}

// GetRepoName returns the test repository name
func GetRepoName() string {
	return Config.RepoName
}

// GetCloudAccountId returns the test cloud account ID
func GetCloudAccountId() string {
	return Config.CloudAccountId
}

// GetOrgId returns the test organization ID
func GetOrgId() string {
	return Config.OrgId
}

// GetControlPolicyId returns the test control policy ID
func GetControlPolicyId() string {
	return Config.ControlPolicyId
}

// GetControlPolicyGroupId returns the test control policy group ID
func GetControlPolicyGroupId() string {
	return Config.ControlPolicyGroupId
}

// GetSlackWebhookUrl returns the test Slack webhook URL
func GetSlackWebhookUrl() string {
	return Config.SlackWebhookUrl
}

// GetSlackAppId returns the test Slack app ID
func GetSlackAppId() string {
	return Config.SlackAppId
}

// GetNotificationEndpointEmail1 returns the first test notification endpoint email
func GetNotificationEndpointEmail1() string {
	return Config.NotificationEndpointEmail1
}

// GetNotificationEndpointEmail2 returns the second test notification endpoint email
func GetNotificationEndpointEmail2() string {
	return Config.NotificationEndpointEmail2
}
