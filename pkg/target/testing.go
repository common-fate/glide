package target

// TestAccessRule returns an AccessRule fixture to be used in tests.
func TestGroup(opt ...func(*Group)) Group {

	ar := Group{
		ID:   "test-target-group",
		Icon: "aws-sso",
		From: From{
			Publisher: "test",
			Name:      "test",
			Version:   "v1.1.1",
		},
	}

	for _, o := range opt {
		o(&ar)
	}

	return ar
}
