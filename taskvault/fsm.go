package taskvault

import (
	"io"

	"github.com/danluki/taskvault/types"
	"github.com/hashicorp/raft"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

// MessageType is the type to encode FSM commands.
type MessageType uint8

const (
	AddPairType MessageType = iota
	DeletePairType
	UpdatePairType
)

// LogApplier is the definition of a function that can apply a Raft log
type LogApplier func(buf []byte, index uint64) interface{}

// LogAppliers is a mapping of the Raft MessageType to the appropriate log
// applier
type LogAppliers map[MessageType]LogApplier

type taskvaultFSM struct {
	store Storage

	// proAppliers holds the set of pro only LogAppliers
	logger *logrus.Entry
}

// NewFSM is used to construct a new FSM with a blank state
func newFSM(store Storage, logger *logrus.Entry) *taskvaultFSM {
	return &taskvaultFSM{
		store:  store,
		logger: logger,
	}
}

// Apply applies a Raft log entry to the key-value store.
func (d *taskvaultFSM) Apply(l *raft.Log) interface{} {
	buf := l.Data
	msgType := MessageType(buf[0])

	d.logger.WithField("command", msgType).Debug("fsm: received command")

	switch msgType {
	case AddPairType:
		return d.applyAddPair(buf[1:])
	case DeletePairType:
		return d.applyDeletePair(buf[1:])
	case UpdatePairType:
		return d.applyUpdatePair(buf[1:])
	}

	return nil
}

func (d *taskvaultFSM) applyAddPair(buf []byte) interface{} {
	var cvr types.CreateValueRequest
	if err := proto.Unmarshal(buf, &cvr); err != nil {
		return err
	}

	err := d.store.SetValue(cvr.Key, cvr.Value)
	if err != nil {
		return err
	}

	return nil
}

func (d *taskvaultFSM) applyDeletePair(buf []byte) interface{} {
	var dpr types.DeleteValueRequest

	if err := proto.Unmarshal(buf, &dpr); err != nil {
		return err
	}

	err := d.store.DeleteValue(dpr.Key)
	if err != nil {
		return err
	}

	return nil
}

func (d *taskvaultFSM) applyUpdatePair(buf []byte) interface{} {
	var uvr types.UpdateValueRequest
	if err := proto.Unmarshal(buf, &uvr); err != nil {
		return err
	}

	err := d.store.UpdateValue(uvr.Key, uvr.Value)
	if err != nil {
		return err
	}

	return nil
}

// Snapshot returns a snapshot of the key-value store. We wrap
// the things we need in taskvaultSnapshot and then send that over to Persist.
// Persist encodes the needed data from taskvaultSnapshot and transport it to
// Restore where the necessary data is replicated into the finite state machine.
// This allows the consensus algorithm to truncate the replicated log.
func (d *taskvaultFSM) Snapshot() (raft.FSMSnapshot, error) {
	return &taskvaultSnapshot{store: d.store}, nil
}

// Restore stores the key-value store to a previous state.
func (d *taskvaultFSM) Restore(r io.ReadCloser) error {
	defer r.Close()
	return d.store.Restore(r)
}

type taskvaultSnapshot struct {
	store Storage
}

func (d *taskvaultSnapshot) Persist(sink raft.SnapshotSink) error {
	if err := d.store.Snapshot(sink); err != nil {
		_ = sink.Cancel()
		return err
	}

	// Close the sink.
	if err := sink.Close(); err != nil {
		return err
	}

	return nil
}

func (d *taskvaultSnapshot) Release() {}
