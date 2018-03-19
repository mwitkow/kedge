package k8sresolver

import (
	"context"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/naming"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)

var (
	testAddr1 = v1.EndpointSubset{
		Ports: []v1.EndpointPort{
			{
				Port: 8080,
				Name: "someName",
			},
			{
				Port: 8081,
				Name: "someName1",
			},
			{
				Port: 8082,
				Name: "someName2",
			},
		},
		Addresses: []v1.EndpointAddress{
			{
				IP: "1.2.3.4",
				TargetRef: &v1.ObjectReference{
					Name:      "pod1",
					Namespace: "ns1",
				},
			},
		},
	}
	modifiedAddr1 = v1.EndpointSubset{
		Ports: []v1.EndpointPort{
			{
				Port: 8080,
				Name: "someName",
			},
		},
		Addresses: []v1.EndpointAddress{
			{
				IP: "1.2.3.5",
				TargetRef: &v1.ObjectReference{
					Name:      "pod1",
					Namespace: "ns1",
				},
			},
		},
	}
	multipleAddrSubset = v1.EndpointSubset{
		Ports: []v1.EndpointPort{
			{
				Port: 8080,
				Name: "someName",
			},
		},
		Addresses: []v1.EndpointAddress{
			{
				IP: "1.2.3.3",
				TargetRef: &v1.ObjectReference{
					Name:      "pod1",
					Namespace: "ns1",
				},
			},
			{
				IP: "1.2.4.4",
				TargetRef: &v1.ObjectReference{
					Name:      "pod2",
					Namespace: "ns1",
				},
			},
			{
				IP: "1.2.5.5",
				TargetRef: &v1.ObjectReference{
					Name:      "pod3",
					Namespace: "ns1",
				},
			},
		},
	}
	someAddrSubset = v1.EndpointSubset{
		Ports: []v1.EndpointPort{
			{
				Port: 8080,
				Name: "someName",
			},
		},
		Addresses: []v1.EndpointAddress{
			{
				IP: "1.2.3.3",
				TargetRef: &v1.ObjectReference{
					Name:      "pod1",
					Namespace: "ns1",
				},
			},
			{
				IP: "1.2.4.4",
				TargetRef: &v1.ObjectReference{
					Name:      "pod2",
					Namespace: "ns1",
				},
			},
		},
	}
	modifiedLastAddrSubset = v1.EndpointSubset{
		Ports: []v1.EndpointPort{
			{
				Port: 8080,
				Name: "someName",
			},
		},
		Addresses: []v1.EndpointAddress{
			{
				IP: "1.2.5.7",
				TargetRef: &v1.ObjectReference{
					Name:      "pod3",
					Namespace: "ns1",
				},
			},
		},
	}
)

func TestWatcher_Next_OK(t *testing.T) {
	for _, tcase := range []struct {
		watchedTargetPort targetPort
		changes           []change
		err               error
		expectedUpdates   [][]*naming.Update
		expectedErrs      []error
	}{
		{
			watchedTargetPort: targetPort{},
			changes:           []change{newTestChange(watch.Added, testAddr1)},
			expectedUpdates: [][]*naming.Update{
				{
					{
						Addr: "1.2.3.4:8080",
						Op:   naming.Add,
					},
				},
			},
		},
		{
			watchedTargetPort: targetPort{isNamed: true, value: "someName2"},
			changes:           []change{newTestChange(watch.Added, testAddr1)},
			expectedUpdates: [][]*naming.Update{
				{
					{
						Addr: "1.2.3.4:8082",
						Op:   naming.Add,
					},
				},
			},
		},
		{
			watchedTargetPort: targetPort{value: "9090"},
			changes:           []change{newTestChange(watch.Added, testAddr1)},
			expectedUpdates: [][]*naming.Update{
				{
					{
						Addr: "1.2.3.4:9090",
						Op:   naming.Add,
					},
				},
			},
		},
		{
			// Non existing port just return no IPs. This makes configuration bit harder to debug, but we cannot assume
			// port is always in any subset.
			watchedTargetPort: targetPort{isNamed: true, value: "non-existing-port-name"},
			changes:           []change{newTestChange(watch.Added, testAddr1)},
			expectedUpdates:   [][]*naming.Update{nil},
		},
		{
			changes: []change{
				newTestChange(watch.Added, testAddr1),
				newTestChange(watch.Modified, modifiedAddr1),
			},
			expectedUpdates: [][]*naming.Update{
				{
					{
						Addr: "1.2.3.4:8080",
						Op:   naming.Add,
					},
				},
				{
					{
						Addr: "1.2.3.4:8080",
						Op:   naming.Delete,
					},
					{
						Addr: "1.2.3.5:8080",
						Op:   naming.Add,
					},
				},
			},
		},
		{
			changes: []change{
				newTestChange(watch.Added, multipleAddrSubset),
				newTestChange(watch.Deleted, someAddrSubset),
				newTestChange(watch.Modified, modifiedLastAddrSubset),
				newTestChange(watch.Deleted, modifiedLastAddrSubset),
			},
			expectedUpdates: [][]*naming.Update{
				{
					{
						Addr: "1.2.3.3:8080",
						Op:   naming.Add,
					},
					{
						Addr: "1.2.4.4:8080",
						Op:   naming.Add,
					},
					{
						Addr: "1.2.5.5:8080",
						Op:   naming.Add,
					},
				},
				{
					{
						Addr: "1.2.3.3:8080",
						Op:   naming.Delete,
					},
					{
						Addr: "1.2.4.4:8080",
						Op:   naming.Delete,
					},
				},
				{
					{
						Addr: "1.2.5.5:8080",
						Op:   naming.Delete,
					},
					{
						Addr: "1.2.5.7:8080",
						Op:   naming.Add,
					},
				},
				{
					{
						Addr: "1.2.5.7:8080",
						Op:   naming.Delete,
					},
				},
			},
		},
		{
			// Test case that did not work previously because of bug.
			changes: []change{
				newTestChange(watch.Added, testAddr1),
				newTestChange(watch.Deleted, testAddr1),
				newTestChange(watch.Added, testAddr1),
			},
			expectedUpdates: [][]*naming.Update{
				{
					{
						Addr: "1.2.3.4:8080",
						Op:   naming.Add,
					},
				},
				{
					{
						Addr: "1.2.3.4:8080",
						Op:   naming.Delete,
					},
				},
				{
					{
						Addr: "1.2.3.4:8080",
						Op:   naming.Add,
					},
				},
			},
		},
		// Malformed state cases. We assume this order of events will never happen:
		{
			changes: []change{
				newTestChange(watch.Deleted, testAddr1),
			},
			expectedErrs: []error{errors.New("malformed internal state for endpoints. On delete event type, we got update for {ns1 pod1} that does not exists in map[]. Doing resync...")},
		},
		{
			changes: []change{
				newTestChange(watch.Modified, testAddr1),
			},
			expectedErrs: []error{errors.New("malformed internal state for endpoints. On modified event type, we got update for {ns1 pod1} that does not exists in map[]. Doing resync...")},
		},
		{
			changes: []change{
				newTestChange(watch.Added, testAddr1),
				newTestChange(watch.Added, testAddr1),
			},
			expectedUpdates: [][]*naming.Update{
				{
					{
						Addr: "1.2.3.4:8080",
						Op:   naming.Add,
					},
				},
			},
			expectedErrs: []error{nil, errors.New("malformed internal state for endpoints. On added event type, we got update for {ns1 pod1} that already exists in map[{ns1 pod1}:1.2.3.4:8080]. Doing resync...")},
		},
	} {
		ok := t.Run("", func(t *testing.T) {
			defer leaktest.CheckTimeout(t, 10*time.Second)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			changeCh := make(chan change, 1)
			errCh := make(chan error, 1)
			w := &watcher{
				ctx:    ctx,
				cancel: cancel,
				target: targetEntry{port: tcase.watchedTargetPort},
				streamer: &streamer{
					changeCh: changeCh,
					errCh:    errCh,
				},
				endpoints: map[key]string{},
			}
			defer w.Close()

			for i, change := range tcase.changes {
				changeCh <- change
				u, err := w.next()
				if len(tcase.expectedErrs) > i && tcase.expectedErrs[i] != nil {
					require.Error(t, err)
					require.Equal(t, tcase.expectedErrs[i].Error(), err.Error())
					continue
				}
				require.NoError(t, err)
				require.Equal(t, tcase.expectedUpdates[i], u, "case %d is wrong", i)
			}
		})
		if !ok {
			return
		}
	}
}
