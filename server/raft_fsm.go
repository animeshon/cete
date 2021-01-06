package server

import (
	"context"
	"errors"
	"io"
	"os"
	"sync"
	"time"

	ceteErrors "github.com/animeshon/cete/errors"
	"github.com/animeshon/cete/marshaler"
	"github.com/animeshon/cete/protobuf"
	"github.com/animeshon/cete/storage"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/raft"
	"go.uber.org/zap"
)

type RaftFSM struct {
	logger *zap.Logger

	kvs        *storage.KVS
	metadata   map[string]*protobuf.Metadata
	nodesMutex sync.RWMutex

	applyCh chan *protobuf.Event
}

func NewRaftFSM(path string, logger *zap.Logger) (*RaftFSM, error) {
	err := os.MkdirAll(path, 0755)
	if err != nil && !os.IsExist(err) {
		logger.Error("failed to make directories", zap.String("path", path), zap.Error(err))
		return nil, err
	}

	kvs, err := storage.NewKVS(path, path, logger)
	if err != nil {
		logger.Error("failed to create key value store", zap.String("path", path), zap.Error(err))
		return nil, err
	}

	// Make sure to garbage collect everything before running the raft (enable with --force-gc-on-startup).
	kvs.RunGC(context.Background(), 0.001)

	// TODO: Context should be passed down to allow for cascade cancellation.
	// TODO: GC should have its own flags for both the interval (--gc-interval=5m) and ratio (--gc-discard-ratio=0.5).
	kvs.ScheduleGC(context.Background(), 5*time.Minute, 0.001)

	return &RaftFSM{
		logger:   logger,
		kvs:      kvs,
		metadata: make(map[string]*protobuf.Metadata, 0),
		applyCh:  make(chan *protobuf.Event, 1024),
	}, nil
}

func (f *RaftFSM) Close() error {
	f.applyCh <- nil
	f.logger.Info("apply channel has closed")

	err := f.kvs.Close()
	if err != nil {
		f.logger.Error("failed to close key value store", zap.Error(err))
		return err
	}
	f.logger.Info("KVS has closed")

	return nil
}

func (f *RaftFSM) Get(key string) ([]byte, error) {
	value, err := f.kvs.Get(key)
	if err != nil {
		if err == ceteErrors.ErrNotFound {
			return nil, err
		}

		f.logger.Error("failed to get value", zap.String("key", key), zap.Error(err))
		return nil, err
	}

	return value, nil
}

func (f *RaftFSM) List(request *protobuf.ListRequest, stream protobuf.KVS_ListServer) error {
	if err := f.kvs.List(request, stream); err != nil {
		f.logger.Error("failed to scan values", zap.String("prefix", request.Prefix), zap.Error(err))
		return err
	}

	return nil
}

func (f *RaftFSM) applySet(key string, value []byte) interface{} {
	err := f.kvs.Set(key, value)
	if err != nil {
		f.logger.Error("failed to set value", zap.String("key", key), zap.Error(err))
		return err
	}

	return nil
}

func (f *RaftFSM) applySetObject(item, meta *protobuf.KeyValuePair, ifMatch, ifNoneMatch string, ifModifiedSince, ifUnmodifiedSince int64) interface{} {
	err := f.kvs.SetObject(item, meta, ifMatch, ifNoneMatch, ifModifiedSince, ifUnmodifiedSince)
	if err != nil {
		f.logger.Debug("failed to set value", zap.String("item.key", item.Key), zap.String("meta.key", meta.Key), zap.Error(err))
		return err
	}

	return nil
}

func (f *RaftFSM) applyDelete(key string) interface{} {
	err := f.kvs.Delete(key)
	if err != nil {
		f.logger.Error("failed to delete value", zap.String("key", key), zap.Error(err))
		return err
	}

	return nil
}

func (f *RaftFSM) applyDeleteObject(itemKey, metaKey, ifMatch, ifNoneMatch string, ifModifiedSince, ifUnmodifiedSince int64) interface{} {
	err := f.kvs.DeleteObject(itemKey, metaKey, ifMatch, ifNoneMatch, ifModifiedSince, ifUnmodifiedSince)
	if err != nil {
		f.logger.Debug("failed to delete value", zap.String("item.key", itemKey), zap.String("meta.key", metaKey), zap.Error(err))
		return err
	}

	return nil
}

func (f *RaftFSM) getMetadata(id string) *protobuf.Metadata {
	if metadata, exists := f.metadata[id]; exists {
		return metadata
	} else {
		f.logger.Warn("metadata not found", zap.String("id", id))
		return nil
	}
}

func (f *RaftFSM) setMetadata(id string, metadata *protobuf.Metadata) {
	f.nodesMutex.Lock()
	f.metadata[id] = metadata
	f.nodesMutex.Unlock()
}

func (f *RaftFSM) deleteMetadata(id string) {
	f.nodesMutex.Lock()
	if _, exists := f.metadata[id]; exists {
		delete(f.metadata, id)
	} else {
		f.logger.Warn("metadata not found", zap.String("id", id))
	}
	f.nodesMutex.Unlock()
}

func (f *RaftFSM) applySetMetadata(id string, metadata *protobuf.Metadata) interface{} {
	f.logger.Debug("set metadata", zap.String("id", id), zap.Any("metadata", metadata))
	f.setMetadata(id, metadata)

	return nil
}

func (f *RaftFSM) applyDeleteMetadata(nodeId string) interface{} {
	f.deleteMetadata(nodeId)

	return nil
}

func (f *RaftFSM) Apply(l *raft.Log) interface{} {
	var event protobuf.Event
	err := proto.Unmarshal(l.Data, &event)
	if err != nil {
		f.logger.Error("failed to unmarshal message bytes to KVS command", zap.Error(err))
		return err
	}

	switch event.Type {
	case protobuf.Event_Join:
		data, err := marshaler.MarshalAny(event.Data)
		if err != nil {
			f.logger.Error("failed to marshal to request from KVS command request", zap.String("type", event.Type.String()), zap.Error(err))
			return err
		}
		if data == nil {
			err = errors.New("nil")
			f.logger.Error("request is nil", zap.String("type", event.Type.String()), zap.Error(err))
			return err
		}
		req := data.(*protobuf.SetMetadataRequest)

		ret := f.applySetMetadata(req.Id, req.Metadata)
		if ret == nil {
			f.applyCh <- &event
		}

		return ret
	case protobuf.Event_Leave:
		data, err := marshaler.MarshalAny(event.Data)
		if err != nil {
			f.logger.Error("failed to marshal to request from KVS command request", zap.String("type", event.Type.String()), zap.Error(err))
			return err
		}
		if data == nil {
			err = errors.New("nil")
			f.logger.Error("request is nil", zap.String("type", event.Type.String()), zap.Error(err))
			return err
		}
		req := *data.(*protobuf.DeleteMetadataRequest)

		ret := f.applyDeleteMetadata(req.Id)
		if ret == nil {
			f.applyCh <- &event
		}

		return ret
	case protobuf.Event_Set:
		data, err := marshaler.MarshalAny(event.Data)
		if err != nil {
			f.logger.Error("failed to marshal to request from KVS command request", zap.String("type", event.Type.String()), zap.Error(err))
			return err
		}
		if data == nil {
			err = errors.New("nil")
			f.logger.Error("request is nil", zap.String("type", event.Type.String()), zap.Error(err))
			return err
		}
		req := *data.(*protobuf.SetRequest)

		ret := f.applySet(req.Key, req.Value)
		if ret == nil {
			f.applyCh <- &event
		}

		return ret
	case protobuf.Event_SetObject:
		data, err := marshaler.MarshalAny(event.Data)
		if err != nil {
			f.logger.Error("failed to marshal to request from KVS command request", zap.String("type", event.Type.String()), zap.Error(err))
			return err
		}
		if data == nil {
			err = errors.New("nil")
			f.logger.Error("request is nil", zap.String("type", event.Type.String()), zap.Error(err))
			return err
		}
		req := *data.(*protobuf.SetObjectRequest)

		ret := f.applySetObject(req.Item, req.Meta, req.IfMatch, req.IfNoneMatch, req.IfModifiedSince, req.IfUnmodifiedSince)
		if ret == nil {
			f.applyCh <- &event
		}

		return ret
	case protobuf.Event_Delete:
		data, err := marshaler.MarshalAny(event.Data)
		if err != nil {
			f.logger.Error("failed to marshal to request from KVS command request", zap.String("type", event.Type.String()), zap.Error(err))
			return err
		}
		if data == nil {
			err = errors.New("nil")
			f.logger.Error("request is nil", zap.String("type", event.Type.String()), zap.Error(err))
			return err
		}
		req := *data.(*protobuf.DeleteRequest)

		ret := f.applyDelete(req.Key)
		if ret == nil {
			f.applyCh <- &event
		}

		return ret
	case protobuf.Event_DeleteObject:
		data, err := marshaler.MarshalAny(event.Data)
		if err != nil {
			f.logger.Error("failed to marshal to request from KVS command request", zap.String("type", event.Type.String()), zap.Error(err))
			return err
		}
		if data == nil {
			err = errors.New("nil")
			f.logger.Error("request is nil", zap.String("type", event.Type.String()), zap.Error(err))
			return err
		}
		req := *data.(*protobuf.DeleteObjectRequest)

		ret := f.applyDeleteObject(req.ItemKey, req.MetaKey, req.IfMatch, req.IfNoneMatch, req.IfModifiedSince, req.IfUnmodifiedSince)
		if ret == nil {
			f.applyCh <- &event
		}

		return ret
	default:
		err = errors.New("command type not support")
		f.logger.Error("unsupported command", zap.String("type", event.Type.String()), zap.Error(err))
		return err
	}
}

func (f *RaftFSM) Stats() map[string]string {
	return f.kvs.Stats()
}

func (f *RaftFSM) Snapshot() (raft.FSMSnapshot, error) {
	return &KVSFSMSnapshot{
		kvs:    f.kvs,
		logger: f.logger,
	}, nil
}

func (f *RaftFSM) Restore(rc io.ReadCloser) error {
	return nil
}

// ---------------------

type KVSFSMSnapshot struct {
	kvs    *storage.KVS
	logger *zap.Logger
}

func (f *KVSFSMSnapshot) Persist(sink raft.SnapshotSink) error {
	start := time.Now()

	f.logger.Info("start to persist items")

	defer func() {
		if err := sink.Close(); err != nil {
			f.logger.Error("failed to close sink", zap.Error(err))
		}
	}()

	if err := f.kvs.Backup(sink); err != nil {
		f.logger.Error("failed to persist items", zap.Error(err))
		return err
	}

	f.logger.Info("finished to persist items", zap.Float64("time", float64(time.Since(start))/float64(time.Second)))

	return nil
}

func (f *KVSFSMSnapshot) Release() {
	f.logger.Info("release")
}
