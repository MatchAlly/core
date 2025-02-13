package cache

func denylistTokenKey(tokenID string) string {
	return "denylist:" + tokenID
}
