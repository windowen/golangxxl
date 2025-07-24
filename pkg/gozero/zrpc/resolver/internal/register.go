package internal

// RegisterResolver registers the direct, etcd and discov schemes to the resolver.
func RegisterResolver() {
	register()
}
