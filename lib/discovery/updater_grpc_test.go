package discovery

import (
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	pb_config "github.com/improbable-eng/kedge/protogen/kedge/config"
	pb_resolvers "github.com/improbable-eng/kedge/protogen/kedge/config/common/resolvers"
	pb_grpcbackends "github.com/improbable-eng/kedge/protogen/kedge/config/grpc/backends"
	pb_grpcroutes "github.com/improbable-eng/kedge/protogen/kedge/config/grpc/routes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdater_OnEvent_AdditionAndModifyAndDelete_GRPC(t *testing.T) {
	defer leaktest.CheckTimeout(t, 10*time.Second)()

	updater := newUpdater(
		&pb_config.DirectorConfig{
			Grpc: &pb_config.DirectorConfig_Grpc{
				Routes: []*pb_grpcroutes.Route{
					{
						Autogenerated:    false,
						AuthorityMatcher: "something",
						PortMatcher:      1234,
						BackendName:      "already_there",
					},
				},
			},
			Http: &pb_config.DirectorConfig_Http{},
		},
		&pb_config.BackendPoolConfig{
			Grpc: &pb_config.BackendPoolConfig_Grpc{
				Backends: []*pb_grpcbackends.Backend{
					{
						Name: "something",
						Resolver: &pb_grpcbackends.Backend_K8S{
							K8S: &pb_resolvers.K8SResolver{
								DnsPortName: "s2.ns1:some-port",
							},
						},
					},
				},
			},
			Http: &pb_config.BackendPoolConfig_Http{},
		},
		"external.example.com",
		"kedge.com/",
	)

	okEvent := event{
		Type: added,
		Object: service{
			Metadata: metadata{
				Name:      "s2",
				Namespace: "ns1",
				Annotations: map[string]string{
					"kedge.com/host-matcher": "s2.external2.host.com",
				},
			},
			Spec: serviceSpec{
				Ports: []portSpec{
					{
						Name:       "grpc",
						Port:       1,
						TargetPort: 1,
					},
					{
						Name:       "grpc-2",
						Port:       2,
						TargetPort: "two",
					},
					{
						Name:       "grpc-2-duplicate",
						Port:       222, // Different matcher.
						TargetPort: "two",
					},
					{
						Name: "grpc-3",
						Port: 3,
					},
					{
						Name:       "not-supported",
						Port:       4,
						TargetPort: 4,
					},
					{
						Name:       "grpctls",
						Port:       447,
						TargetPort: "six",
					},
				},
			},
		},
	}

	d, b, err := updater.onEvent(okEvent)
	require.NoError(t, err)

	expectedDirectorConfig := &pb_config.DirectorConfig{
		Grpc: &pb_config.DirectorConfig_Grpc{
			Routes: []*pb_grpcroutes.Route{
				{
					Autogenerated:    false,
					AuthorityMatcher: "something",
					PortMatcher:      1234,
					BackendName:      "already_there",
				},
				{
					Autogenerated:    true,
					AuthorityMatcher: "s2.external2.host.com",
					PortMatcher:      1,
					BackendName:      "s2_ns1_1",
				},
				{
					Autogenerated:    true,
					AuthorityMatcher: "s2.external2.host.com",
					PortMatcher:      2,
					BackendName:      "s2_ns1_two",
				},
				{
					Autogenerated:    true,
					AuthorityMatcher: "s2.external2.host.com",
					PortMatcher:      3,
					BackendName:      "s2_ns1_3",
				},
				{
					Autogenerated:    true,
					AuthorityMatcher: "s2.external2.host.com",
					PortMatcher:      222,
					BackendName:      "s2_ns1_two",
				},
				{
					Autogenerated:    true,
					AuthorityMatcher: "s2.external2.host.com",
					PortMatcher:      447,
					BackendName:      "s2_ns1_six",
				},
			},
		},
		Http: &pb_config.DirectorConfig_Http{},
	}
	assert.Equal(t, expectedDirectorConfig, d)

	expectedBackendpoolConfig := &pb_config.BackendPoolConfig{
		Grpc: &pb_config.BackendPoolConfig_Grpc{
			Backends: []*pb_grpcbackends.Backend{
				{
					Autogenerated: true,
					Name:          "s2_ns1_1",
					Resolver: &pb_grpcbackends.Backend_K8S{
						K8S: &pb_resolvers.K8SResolver{
							DnsPortName: "s2.ns1:grpc",
						},
					},
				},
				{
					Autogenerated: true,
					Name:          "s2_ns1_3",
					Resolver: &pb_grpcbackends.Backend_K8S{
						K8S: &pb_resolvers.K8SResolver{
							DnsPortName: "s2.ns1:grpc-3",
						},
					},
				},
				{
					Autogenerated: true,
					Name:          "s2_ns1_six",
					Resolver: &pb_grpcbackends.Backend_K8S{
						K8S: &pb_resolvers.K8SResolver{
							DnsPortName: "s2.ns1:grpctls",
						},
					},
					Security: &pb_grpcbackends.Security{
						InsecureSkipVerify: true,
					},
				},
				{
					Autogenerated: true,
					Name:          "s2_ns1_two",
					Resolver: &pb_grpcbackends.Backend_K8S{
						K8S: &pb_resolvers.K8SResolver{
							DnsPortName: "s2.ns1:grpc-2-duplicate",
						},
					},
				},
				{
					Name: "something",
					Resolver: &pb_grpcbackends.Backend_K8S{
						K8S: &pb_resolvers.K8SResolver{
							DnsPortName: "s2.ns1:some-port",
						},
					},
				},
			},
		},
		Http: &pb_config.BackendPoolConfig_Http{},
	}
	assert.Equal(t, expectedBackendpoolConfig, b)

	modifyEvent := event{
		Type: modified,
		Object: service{
			Kind: "Services",
			Metadata: metadata{
				Name:      "s2",
				Namespace: "ns1",
			},
			Spec: serviceSpec{
				Ports: []portSpec{
					{
						Name:       "grpc",
						Port:       1,
						TargetPort: 1,
					},
					{
						Name:       "grpc-4",
						Port:       4,
						TargetPort: "four",
					},
					{
						Name: "grpc-5",
						Port: 5,
					},
					{
						Name:       "not-supported",
						Port:       4,
						TargetPort: 4,
					},
				},
			},
		},
	}

	d, b, err = updater.onEvent(modifyEvent)
	require.NoError(t, err)

	expectedDirectorConfig2 := &pb_config.DirectorConfig{
		Grpc: &pb_config.DirectorConfig_Grpc{
			Routes: []*pb_grpcroutes.Route{
				{
					Autogenerated:    false,
					AuthorityMatcher: "something",
					PortMatcher:      1234,
					BackendName:      "already_there",
				},
				{
					Autogenerated:    true,
					AuthorityMatcher: "s2.external.example.com",
					PortMatcher:      1,
					BackendName:      "s2_ns1_1",
				},
				{
					Autogenerated:    true,
					AuthorityMatcher: "s2.external.example.com",
					PortMatcher:      4,
					BackendName:      "s2_ns1_four",
				},
				{
					Autogenerated:    true,
					AuthorityMatcher: "s2.external.example.com",
					PortMatcher:      5,
					BackendName:      "s2_ns1_5",
				},
			},
		},
		Http: &pb_config.DirectorConfig_Http{},
	}
	assert.Equal(t, expectedDirectorConfig2, d)

	expectedBackendpoolConfig2 := &pb_config.BackendPoolConfig{
		Grpc: &pb_config.BackendPoolConfig_Grpc{
			Backends: []*pb_grpcbackends.Backend{
				{
					Autogenerated: true,
					Name:          "s2_ns1_1",
					Resolver: &pb_grpcbackends.Backend_K8S{
						K8S: &pb_resolvers.K8SResolver{
							DnsPortName: "s2.ns1:grpc",
						},
					},
				},
				{
					Autogenerated: true,
					Name:          "s2_ns1_5",
					Resolver: &pb_grpcbackends.Backend_K8S{
						K8S: &pb_resolvers.K8SResolver{
							DnsPortName: "s2.ns1:grpc-5",
						},
					},
				},
				{
					Autogenerated: true,
					Name:          "s2_ns1_four",
					Resolver: &pb_grpcbackends.Backend_K8S{
						K8S: &pb_resolvers.K8SResolver{
							DnsPortName: "s2.ns1:grpc-4",
						},
					},
				},
				{
					Name: "something",
					Resolver: &pb_grpcbackends.Backend_K8S{
						K8S: &pb_resolvers.K8SResolver{
							DnsPortName: "s2.ns1:some-port",
						},
					},
				},
			},
		},
		Http: &pb_config.BackendPoolConfig_Http{},
	}
	assert.Equal(t, expectedBackendpoolConfig2, b)

	deleteEvent := event{
		Type: deleted,
		Object: service{
			Kind: "Services",
			Metadata: metadata{
				Name:      "s2",
				Namespace: "ns1",
			},
		},
	}

	d, b, err = updater.onEvent(deleteEvent)
	require.NoError(t, err)

	expectedDirectorConfig3 := &pb_config.DirectorConfig{
		Grpc: &pb_config.DirectorConfig_Grpc{
			Routes: []*pb_grpcroutes.Route{
				{
					Autogenerated:    false,
					AuthorityMatcher: "something",
					PortMatcher:      1234,
					BackendName:      "already_there",
				},
			},
		},
		Http: &pb_config.DirectorConfig_Http{},
	}
	assert.Equal(t, expectedDirectorConfig3, d)

	expectedBackendpoolConfig3 := &pb_config.BackendPoolConfig{
		Grpc: &pb_config.BackendPoolConfig_Grpc{
			Backends: []*pb_grpcbackends.Backend{
				{
					Name: "something",
					Resolver: &pb_grpcbackends.Backend_K8S{
						K8S: &pb_resolvers.K8SResolver{
							DnsPortName: "s2.ns1:some-port",
						},
					},
				},
			},
		},
		Http: &pb_config.BackendPoolConfig_Http{},
	}
	assert.Equal(t, expectedBackendpoolConfig3, b)
}
