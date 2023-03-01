resource "commonfate_access_rule" "test" {
  name ="test"
  description="test"
  groups=[
        "common_fate_administrators",
  ]
  target=[
   
	{
        field="accountId",
				value="632700053629"
	},
	
	{
        field="permissionSetArn",
				value="arn:aws:sso:::permissionSet/ssoins-825968feece9a0b6/ps-dda57372ebbfeb94"
	},
	
	
  ]
  target_provider_id="aws-sso-v2"
  duration=3600
}
resource "commonfate_access_rule" "asdflkjasdf" {
  name ="asdflkjasdf"
  description="f;aslkfasjdk;f"
  groups=[
        "common_fate_administrators",
        "granted_administrators",
  ]
  target=[
   
	{
        field="permissionSetArn",
				value="arn:aws:sso:::permissionSet/ssoins-825968feece9a0b6/ps-dda57372ebbfeb94"
	},
	
	
	{
        field="accountId",
				value=["234148001710","829873318951","632700053629",]
	},
	
  ]
  target_provider_id="aws-sso-v2"
  duration=3600
}
