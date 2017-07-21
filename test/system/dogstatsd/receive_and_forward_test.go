package dogstatsd_test

import (
	"encoding/json"
	"testing"
	"time"

	log "github.com/cihub/seelog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DataDog/datadog-agent/pkg/metrics"
	"github.com/DataDog/datadog-agent/pkg/metadata/v5"
)

func testMetadata(t *testing.T, d *dogstatsdTest) {
	// waiting for metadata payload
	timeOut := time.Tick(10 * time.Second)
	select {
	case <-d.requestReady:
	case <-timeOut:
		require.Fail(t, "Timeout: the backend never receive a metadata requests from dogstatsd")
	}

	requests := d.getRequests()
	require.Len(t, requests, 1)

	metadata := v5.Payload{}
	err := json.Unmarshal([]byte(requests[0]), &metadata)
	require.NoError(t, err, "Could not Unmarshal metadata request")

	require.NotNil(t, metadata.Os)
}

func TestReceiveAndForward(t *testing.T) {
	d := setupDogstatsd(t)
	defer d.teardown()
	defer log.Flush()

	testMetadata(t, d)

	d.sendUDP("_sc|test.ServiceCheck|0")

	timeOut := time.Tick(30 * time.Second)
	select {
	case <-d.requestReady:
	case <-timeOut:
		require.Fail(t, "Timeout: the backend never receive a requests from dogstatsd")
	}

	requests := d.getRequests()
	require.Len(t, requests, 1)

	sc := []metrics.ServiceCheck{}
	err := json.Unmarshal([]byte(requests[0]), &sc)
	require.NoError(t, err, "Could not Unmarshal request")

	require.Len(t, sc, 2)
	assert.Equal(t, sc[0].CheckName, "test.ServiceCheck")
	assert.Equal(t, sc[0].Status, metrics.ServiceCheckOK)

	assert.Equal(t, sc[1].CheckName, "datadog.agent.up")
	assert.Equal(t, sc[1].Status, metrics.ServiceCheckOK)
}
