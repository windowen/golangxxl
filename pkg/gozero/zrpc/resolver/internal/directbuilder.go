package internal

import (
	"strings"

	"google.golang.org/grpc/resolver"

	"queueJob/pkg/gozero/zrpc/resolver/internal/targets"
)

type directBuilder struct{}

func (d *directBuilder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (
	resolver.Resolver, error) {
	endpoints := strings.FieldsFunc(targets.GetEndpoints(target), func(r rune) bool {
		return r == EndpointSepChar
	})
	endpoints = subset(endpoints, subsetSize)
	addrs := make([]resolver.Address, 0, len(endpoints))

	for _, val := range endpoints {
		addrs = append(addrs, resolver.Address{
			Addr: val,
		})
	}
	if err := cc.UpdateState(resolver.State{
		Addresses: addrs,
	}); err != nil {
		return nil, err
	}

	return &nopResolver{cc: cc}, nil
}

func (d *directBuilder) Scheme() string {
	return DirectScheme
}
