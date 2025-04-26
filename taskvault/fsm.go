package taskvault

import (
	"io"

	"github.com/danluki/taskvault/pkg/types"
	"github.com/hashicorp/raft"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type MessageType uint8

const (
	AddPairType MessageType = iota
	DeletePairType
	UpdatePairType
)

type Pair struct {
	Key   string
	Value string
}

type LogApplier func(buf []byte, index uint64) interface{}

type LogAppliers map[MessageType]LogApplier

type taskvaultFSM struct {
	store SyncraStorage

	logger *zap.SugaredLogger
}

func newFSM(store SyncraStorage, logger *zap.SugaredLogger) *taskvaultFSM {
	return &taskvaultFSM{
		store:  store,
		logger: logger,
	}
}

func (d *taskvaultFSM) Apply(l *raft.Log) interface{} {
	buf := l.Data
	msgType := MessageType(buf[0])

	d.logger.Debug("fsm: received command", zap.Int8("command", int8(msgType)))

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

func (d *taskvaultFSM) Snapshot() (raft.FSMSnapshot, error) {
	return &taskvaultSnapshot{store: d.store}, nil
}

func (d *taskvaultFSM) Restore(r io.ReadCloser) error {
	defer r.Close()
	return d.store.Restore(r)
}

type taskvaultSnapshot struct {
	store SyncraStorage
}

func (d *taskvaultSnapshot) Persist(sink raft.SnapshotSink) error {
	if err := d.store.Snapshot(sink); err != nil {
		_ = sink.Cancel()
		return err
	}

	if err := sink.Close(); err != nil {
		return err
	}

	return nil
}

func (d *taskvaultSnapshot) Release() {}
