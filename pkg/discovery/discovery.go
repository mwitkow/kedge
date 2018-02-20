package discovery

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/improbable-eng/kedge/pkg/k8s"
	"github.com/improbable-eng/kedge/pkg/sharedflags"
	pb_config "github.com/improbable-eng/kedge/protogen/kedge/config"
	"github.com/jpillora/backoff"
	"github.com/mwitkow/go-flagz/protobuf"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	selectorKeySuffix           = "kedge-exposed"
	hostMatcherAnnotationSuffix = "host-matcher"
)

var (
	// TODO(bplotka): Consider moving to regex with .namespace .name variables.
	flagExternalDomainSuffix = sharedflags.Set.String("discovery_external_domain_suffix", "", "Required suffix "+
		"that will be added to service name to constructs external domain for director route")
	flagAnnotationLabelPrefix = sharedflags.Set.String("discovery_label_annotation_prefix", "kedge.com",
		"Annotation/label prefix for all kedge annotations and kedge-exposed label")
	flagDiscoveryByAnnotation = sharedflags.Set.Bool("enable_discovery_by_annotation", false,
		"Filter discovered endpoints by annotation label. Default is to discover all endpoints.")
	streamRetryBackoff = &backoff.Backoff{
		Min:    50 * time.Millisecond,
		Jitter: true,
		Factor: 2,
		Max:    2 * time.Second,
	}
)

// RoutingDiscovery allows to get fresh director and backendpool configuration filled with autogenerated routings based on service annotations.
// See more details about it in server/README.md#Running_with_Dynaming_Routing_Discovery
type RoutingDiscovery struct {
	logger                logrus.FieldLogger
	serviceClient         serviceClient
	baseBackendpool       *pb_config.BackendPoolConfig
	baseDirector          *pb_config.DirectorConfig
	labelSelectorKey      string
	externalDomainSuffix  string
	labelAnnotationPrefix string
}

// NewFromFlags creates new RoutingDiscovery flow flags.
func NewFromFlags(logger logrus.FieldLogger, baseDirector *pb_config.DirectorConfig, baseBackendpool *pb_config.BackendPoolConfig) (*RoutingDiscovery, error) {
	if *flagExternalDomainSuffix == "" {
		return nil, errors.Errorf("required flag 'discovery_external_domain_suffix' is not specified.")
	}

	apiClient, err := k8s.NewFromFlags()
	if err != nil {
		return nil, err
	}
	return NewWithClient(logger, baseDirector, baseBackendpool, &client{k8sClient: apiClient}), nil
}

// NewWithClient returns a new Kubernetes RoutingDiscovery using given k8s.APIClient configured to be used against kube-apiserver.
func NewWithClient(logger logrus.FieldLogger, baseDirector *pb_config.DirectorConfig, baseBackendpool *pb_config.BackendPoolConfig, serviceClient serviceClient) *RoutingDiscovery {
	// If we want to discover by label selector, set a key.
	var labelSelectorKey string
	if *flagDiscoveryByAnnotation {
		labelSelectorKey = fmt.Sprintf("%s%s", *flagAnnotationLabelPrefix, selectorKeySuffix)
	}

	return &RoutingDiscovery{
		logger:                logger,
		baseBackendpool:       baseBackendpool,
		baseDirector:          baseDirector,
		serviceClient:         serviceClient,
		labelSelectorKey:      labelSelectorKey,
		externalDomainSuffix:  *flagExternalDomainSuffix,
		labelAnnotationPrefix: *flagAnnotationLabelPrefix,
	}
}

// DiscoverOnce returns director & backendpool configs filled with mix of persistent routes & backends given in base configs and dynamically discovered ones.
func (d *RoutingDiscovery) DiscoverOnce(ctx context.Context, timeToWait time.Duration) (*pb_config.DirectorConfig, *pb_config.BackendPoolConfig, error) {
	ctx, cancel := context.WithTimeout(ctx, timeToWait) // Let's give 4 seconds to gather all changes.
	defer cancel()

	watchResultCh := make(chan watchResult)
	defer close(watchResultCh)

	err := startWatchingServicesChanges(ctx, d.labelSelectorKey, d.serviceClient, watchResultCh)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Failed to start watching services by %s selector stream", d.labelSelectorKey)
	}

	updater := newUpdater(
		d.baseDirector,
		d.baseBackendpool,
		d.externalDomainSuffix,
		d.labelAnnotationPrefix,
	)

	var resultDirectorConfig *pb_config.DirectorConfig
	var resultBackendPool *pb_config.BackendPoolConfig
	for {
		var event event
		select {
		case <-ctx.Done():
			if resultBackendPool == nil {
				resultBackendPool = d.baseBackendpool
			}
			if resultDirectorConfig == nil {
				resultDirectorConfig = d.baseDirector
			}
			// Time is up, let's return what we have so far.
			return resultDirectorConfig, resultBackendPool, nil
		case r := <-watchResultCh:
			if r.err != nil {
				return nil, nil, errors.Wrap(r.err, "error on reading event stream")
			}
			event = *r.ep
		}

		resultDirectorConfig, resultBackendPool, err = updater.onEvent(event)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "error on updating routing on event %v", event)
		}
	}
}

// DiscoverAndSetFlags constantly watches service list endpoint for changes and sets director and backendpool
// flags to change routing configuration. Having it set by flagz is giving us possibility to see current values in
// debug/flagz page and perform proper apply in different place (where we are parsing flags).
func (d *RoutingDiscovery) DiscoverAndSetFlags(
	ctx context.Context,
	directorFlagz *protoflagz.DynProto3Value,
	backendpoolFlagz *protoflagz.DynProto3Value,
) error {
	watchResultCh := make(chan watchResult)
	defer close(watchResultCh)

	for ctx.Err() == nil {
		streamCtx, streamCancel := context.WithCancel(ctx)
		// All errors from watching service or watchResultCh are irrecoverable.
		err := startWatchingServicesChanges(streamCtx, d.labelSelectorKey, d.serviceClient, watchResultCh)
		if err != nil {
			streamCancel()
			d.logger.WithError(err).Errorf("discovery: Failed to start watching services by %s selector stream", d.labelSelectorKey)

			time.Sleep(streamRetryBackoff.Duration())
			continue
		}
		streamRetryBackoff.Reset()

		updater := newUpdater(
			d.baseDirector,
			d.baseBackendpool,
			d.externalDomainSuffix,
			d.labelAnnotationPrefix,
		)

		aggregateDelay := 300 * time.Millisecond
		maxAggrUpdates := 50
		for ctx.Err() == nil {
			// Watch loop.

			director, backendPool, err := d.nextWithAggregate(ctx, updater, watchResultCh, aggregateDelay, maxAggrUpdates)
			if err != nil {
				d.logger.WithError(err).Error("discovery: Error on watching event stream. Retrying stream... ")
				break
			}

			d.logger.Infof("Setting director and backendpool configs with %d (httpScheme) and %d (grpcScheme) backends.",
				len(backendPool.Http.Backends), len(backendPool.Grpc.Backends))

			// Only first update should aggregate that much, to not erase and create backend unnecessarily on watch endpoint EOF (which is every 15 minutes).
			aggregateDelay = 100 * time.Millisecond
			maxAggrUpdates = 10

			err = d.setFlags(director, directorFlagz, backendPool, backendpoolFlagz)
			if err != nil {
				streamCancel()
				// This is critical error, since retry will not help. Not much we can do, abort dynamic discovery.
				return errors.Wrap(err, "discovery: Critical error on updating flags from director and backendpool configs. Discover will not work!")
			}
		}
		streamCancel()
	}

	return nil
}

// nextWithAggregate is waiting for next watcher update. It does not immediately return valid full director and
// backendPool configs that includes this single change. Instead it aggregates all changes that are gathered within next "aggregateDelay"
// after last update to aggregate these.
// Parameter called "maxAggrUpdates" is to have safeguard in case of constant changes. In theory time waiting for an aggregated
// changes in our configs is maxAggrUpdates * aggregateDelay if we have constant changes (unlikely).
// NOTE that error from next() method are irrecoverable.
func (d *RoutingDiscovery) nextWithAggregate(ctx context.Context, updater *updater, watchResultCh <-chan watchResult, aggregateDelay time.Duration, maxAggrUpdates int) (*pb_config.DirectorConfig, *pb_config.BackendPoolConfig, error) {
	var firstUpdateAt time.Time

	var director *pb_config.DirectorConfig
	var backendPool *pb_config.BackendPoolConfig
	for {
		var event event
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		case r := <-watchResultCh:
			if r.err != nil {
				return nil, nil, r.err
			}

			if firstUpdateAt.IsZero() {
				firstUpdateAt = time.Now()
			}
			event = *r.ep
		case <-time.After(aggregateDelay):
			if firstUpdateAt.IsZero() {
				// Still waiting for any update.
				continue
			}
			return director, backendPool, nil
		}

		// Safe guard in case of constant changes.
		maxAggrUpdates--
		if maxAggrUpdates <= 0 {
			return director, backendPool, nil
		}

		var err error
		director, backendPool, err = updater.onEvent(event)
		if err != nil {
			// There is possibility we missed event, so retry stream.
			return nil, nil, errors.Wrapf(err, "internal error on updating routing on event %v.", event)
		}
	}
}

func (d *RoutingDiscovery) setFlags(
	directorConfig *pb_config.DirectorConfig,
	directorFlagz *protoflagz.DynProto3Value,
	backenpoolConfig *pb_config.BackendPoolConfig,
	backendpoolFlagz *protoflagz.DynProto3Value,
) error {
	director, err := (&jsonpb.Marshaler{}).MarshalToString(directorConfig)
	if err != nil {
		return errors.Wrap(err, "failed to marshal directorConfig")
	}
	err = directorFlagz.Set(string(director))
	if err != nil {
		return errors.Wrap(err, "failed to set directorConfig into protoflagz flag")
	}

	backendpool, err := (&jsonpb.Marshaler{}).MarshalToString(backenpoolConfig)
	if err != nil {
		return errors.Wrap(err, "failed to marshal backenpoolConfig")
	}
	err = backendpoolFlagz.Set(backendpool)
	if err != nil {
		return errors.Wrap(err, "failed to set backenpoolConfig into protoflagz flag")
	}
	return nil
}
