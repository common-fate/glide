# gon.hcl
#
# The path follows a pattern
# ./dist/BUILD-ID_TARGET/BINARY-NAME
source    = ["./dist/gdeploy-macos-arm_darwin_arm64/gdeploy"]
bundle_id = "io.commonfate.gdeploy"

apple_id {
  username = "chris@exponentlabs.io"
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "Developer ID Application: Common Fate Technologies Pty Ltd"
}
