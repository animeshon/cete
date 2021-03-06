package storage

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	_errors "errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/animeshon/cete/errors"
	"github.com/animeshon/cete/protobuf"
	"github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"
	"github.com/dgraph-io/badger/v3/y"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type KVS struct {
	dir      string
	valueDir string
	db       *badger.DB
	logger   *zap.Logger
}

func NewKVS(dir string, valueDir string, logger *zap.Logger) (*KVS, error) {
	opts := badger.DefaultOptions(dir)
	opts.ValueDir = valueDir
	opts.SyncWrites = !strings.Contains(os.Getenv("FLAGS"), "--disable-sync-writes")
	opts.Logger = NewBadgerLogger(logger)
	// opts.Truncate = strings.Contains(os.Getenv("FLAGS"), "--truncate")
	opts.Compression = options.Snappy
	opts.ValueLogFileSize = 1<<29 - 1

	db, err := badger.Open(opts)
	if err != nil {
		logger.Error("failed to open database", zap.Any("opts", opts), zap.Error(err))
		return nil, err
	}

	return &KVS{
		dir:      dir,
		valueDir: valueDir,
		db:       db,
		logger:   logger,
	}, nil
}

func (k *KVS) Close() error {
	if err := k.db.Close(); err != nil {
		k.logger.Error("failed to close database", zap.Error(err))
		return err
	}

	return nil
}

func (k *KVS) RunGC(ctx context.Context, discardRatio float64) {
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := k.db.RunValueLogGC(discardRatio)
			if err != nil {
				if err == badger.ErrNoRewrite {
					goto finished
				}

				k.logger.Error("garbage collection failed", zap.Error(err))
				goto finished
			}
		}
	}

finished:
	k.logger.Info("garbage collection finished", zap.Float64("time", float64(time.Since(start))/float64(time.Second)))
}

func (k *KVS) ScheduleGC(ctx context.Context, interval time.Duration, discardRatio float64) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				k.RunGC(ctx, discardRatio)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (k *KVS) Get(key string) ([]byte, error) {
	start := time.Now()

	var value []byte
	if err := k.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return err
			}

			k.logger.Error("failed to get item", zap.String("key", key), zap.Error(err))
			return err
		}

		err = item.Value(func(val []byte) error {
			value = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			k.logger.Error("failed to get item value", zap.String("key", key), zap.Error(err))
			return err
		}

		return nil
	}); err == badger.ErrKeyNotFound {
		k.logger.Debug("not found", zap.String("key", key), zap.Error(err))
		return nil, errors.ErrNotFound
	} else if err != nil {
		k.logger.Error("failed to get value", zap.String("key", key), zap.Error(err))
		return nil, err
	}

	k.logger.Debug("get", zap.String("key", key), zap.Float64("time", float64(time.Since(start))/float64(time.Second)))
	return value, nil
}

func (k *KVS) Scan(prefix string) ([][]byte, error) {
	start := time.Now()

	var value [][]byte
	if err := k.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefixBytes := []byte(prefix)
		for it.Seek(prefixBytes); it.ValidForPrefix(prefixBytes); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				value = append(value, append([]byte{}, val...))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		k.logger.Error("failed to scan value", zap.String("prefix", prefix), zap.Error(err))
		return nil, err
	}

	k.logger.Debug("scan", zap.String("prefix", prefix), zap.Float64("time", float64(time.Since(start))/float64(time.Second)))
	return value, nil
}

func (k *KVS) Set(key string, value []byte) error {
	start := time.Now()

	if err := k.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), value)
		if err != nil {
			k.logger.Error("failed to set item", zap.String("key", key), zap.Error(err))
			return err
		}
		return nil
	}); err != nil {
		k.logger.Error("failed to set value", zap.String("key", key), zap.Error(err))
		return err
	}

	k.logger.Debug("set", zap.String("key", key), zap.Float64("time", float64(time.Since(start))/float64(time.Second)))
	return nil
}

func GetObjectMetaVersion1(value []byte) (string, error) {
	type T struct {
		Version string `json:"version"`
	}

	var _meta *T
	if err := json.Unmarshal(value, &_meta); err != nil {
		return "", err
	}

	return _meta.Version, nil
}

func GetObjectMetaVersion2(value []byte) (string, error) {
	r, err := gzip.NewReader(bytes.NewReader(value))
	if err != nil {
		return "", err
	}
	defer r.Close()

	var b bytes.Buffer
	if _, err := io.Copy(&b, r); err != nil {
		return "", err
	}

	var _meta protobuf.ObjectMeta
	if err := proto.Unmarshal(b.Bytes(), &_meta); err != nil {
		return "", err
	}

	return _meta.Version, nil
}

func GetObjectMetaVersion(value []byte) (string, error) {
	if value[0] == '{' && value[len(value)-1] == '}' {
		return GetObjectMetaVersion1(value)
	}

	// The following is the better disk-efficient implementation.
	return GetObjectMetaVersion2(value)
}

func (k *KVS) SetObject(item, meta *protobuf.KeyValuePair, ifMatch, ifNoneMatch string, ifModifiedSince, ifUnmodifiedSince int64) error {
	start := time.Now()

	// TODO: Implement [If-Modified-Since] and [If-Unmodified-Since] logic.
	if err := k.db.Update(func(txn *badger.Txn) error {
		badgerItem, err := txn.Get([]byte(meta.Key))
		if err != nil && !_errors.Is(err, badger.ErrKeyNotFound) {
			k.logger.Error("failed to get item metadata", zap.String("key", meta.Key), zap.Error(err))
			return err
		}

		// If-Match fails if the latest metadata version does not match the provided version or if it doesn't exist.
		if len(ifMatch) != 0 && err == badger.ErrKeyNotFound {
			return _errors.New("precondition failed [If-Match]: metadata does not exist")
		}

		// If-None-Match fails if the latest metadata version matches the provided version.
		// In case If-None-Match is set to '*' it fails if any metadata exists, no matter its version.
		if ifNoneMatch == "*" && err != badger.ErrKeyNotFound {
			return _errors.New("precondition failed [If-None-Match]: metadata already exists")
		}

		// Verify that the metadata version matches the precondition.
		if len(ifMatch) != 0 {
			var value []byte
			err = badgerItem.Value(func(val []byte) error {
				value = val
				return nil
			})
			if err != nil {
				k.logger.Error("failed to get metadata value", zap.String("key", meta.Key), zap.Error(err))
				return err
			}

			if len(value) < 2 {
				return _errors.New("precondition failed [If-Match]: metadata is too short")
			}

			version, err := GetObjectMetaVersion(value)
			if err != nil {
				return err
			}

			if version != ifMatch {
				return _errors.New("precondition failed [If-Match]: value mismatch")
			}
		}

		err = txn.Set([]byte(item.Key), item.Value)
		if err != nil {
			k.logger.Error("failed to set item", zap.String("key", item.Key), zap.Error(err))
			return err
		}

		err = txn.Set([]byte(meta.Key), meta.Value)
		if err != nil {
			k.logger.Error("failed to set item metadata", zap.String("key", meta.Key), zap.Error(err))
			return err
		}
		return nil
	}); err != nil {
		k.logger.Debug("failed to set value", zap.String("key", item.Key), zap.Error(err))
		return err
	}

	k.logger.Debug("set", zap.String("key", item.Key), zap.Float64("time", float64(time.Since(start))/float64(time.Second)))
	return nil
}

func (k *KVS) Delete(key string) error {
	start := time.Now()

	if err := k.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		if err != nil {
			k.logger.Error("failed to delete item", zap.String("key", key), zap.Error(err))
			return err
		}
		return nil
	}); err != nil {
		k.logger.Error("failed to delete value", zap.String("key", key), zap.Error(err))
		return err
	}

	k.logger.Debug("delete", zap.String("key", key), zap.Float64("time", float64(time.Since(start))/float64(time.Second)))
	return nil
}

func (k *KVS) DeleteObject(itemKey, metaKey string, ifMatch, ifNoneMatch string, ifModifiedSince, ifUnmodifiedSince int64) error {
	start := time.Now()

	// TODO: Implement [If-Match], [If-None-Match], [If-Modified-Since] and [If-Unmodified-Since] logic.
	if err := k.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(itemKey))
		if err != nil {
			k.logger.Error("failed to delete item", zap.String("key", itemKey), zap.Error(err))
			return err
		}
		err = txn.Delete([]byte(metaKey))
		if err != nil {
			k.logger.Error("failed to delete item", zap.String("key", metaKey), zap.Error(err))
			return err
		}
		return nil
	}); err != nil {
		k.logger.Debug("failed to delete value", zap.String("key", itemKey), zap.Error(err))
		return err
	}

	k.logger.Debug("delete", zap.String("key", itemKey), zap.Float64("time", float64(time.Since(start))/float64(time.Second)))
	return nil
}

func (k *KVS) Stats() map[string]string {
	stats := map[string]string{}

	stats["num_reads"] = y.NumReads.String()
	stats["num_writes"] = y.NumWrites.String()
	stats["num_bytes_read"] = y.NumBytesRead.String()
	stats["num_bytes_written"] = y.NumBytesWritten.String()
	stats["num_lsm_gets"] = y.NumLSMGets.String()
	stats["num_lsm_bloom_Hits"] = y.NumLSMBloomHits.String()
	stats["num_gets"] = y.NumGets.String()
	stats["num_puts"] = y.NumPuts.String()
	stats["num_blocked_puts"] = y.NumBlockedPuts.String()
	stats["num_memtables_gets"] = y.NumMemtableGets.String()
	stats["lsm_size"] = y.LSMSize.String()
	stats["vlog_size"] = y.VlogSize.String()
	stats["pending_writes"] = y.PendingWrites.String()

	return stats
}

func (k *KVS) Backup(w io.Writer) error {
	start := time.Now()

	// TODO: Make snapshots incremental using this functionality.
	ts, err := k.db.Backup(w, 0)
	if err != nil {
		return err
	}

	k.logger.Info("backup succesfully completed", zap.Uint64("ts", ts), zap.Float64("time", float64(time.Since(start))/float64(time.Second)))
	return nil
}
