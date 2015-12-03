package bulk

import (
	"sync"
	"time"

	"gopkg.in/underarmour/dynago.v1"
)

// Configuration for the bulk writer
type Config struct {
	Client *dynago.Client // The dynago client we're using to make requests
	Table  string         // The table we're writing to.

	// How many goroutines to run that execute write operations. This
	// effectively sets the max parallel writes we can do on the table.
	Concurrency int // Defaults to 1 if unset

	// How many records to write per bulk write.
	PerWrite int // Defaults to 25 if unset
}

func (c *Config) setDefaults() {
	if c.Concurrency < 1 {
		c.Concurrency = 1
	}

	if c.PerWrite < 1 {
		c.PerWrite = 25
	}
}

// Create a new BulkWriter.
func New(config Config) *BulkWriter {
	config.setDefaults()

	writer := &BulkWriter{
		client:  config.Client,
		table:   config.Table,
		ch:      make(chan message, config.Concurrency*10),
		groups:  make(chan group),
		results: make(chan Result),
	}
	writer.wg.Add(1)
	go writer.main(config.PerWrite)
	for i := 0; i < config.Concurrency; i++ {
		writer.wg.Add(1)
		go writer.worker(i)
	}
	return writer
}

/*
BulkWriter manages a series of bulk operations.
*/
type BulkWriter struct {
	client  *dynago.Client
	table   string
	ch      chan message
	groups  chan group
	results chan Result
	wg      sync.WaitGroup
}

/*
Queue up a write.
This function is safe to call from any number of goroutines. It will
block only if all workers are presently executing batches plus our buffer
is full, creating built-in back-pressure on writing.
*/
func (b *BulkWriter) Write(doc dynago.Document) {
	b.ch <- message{doc: doc}
}

/*
Queue up a delete.
*/
func (b *BulkWriter) Delete(key dynago.Document) {
	b.ch <- message{deleteKey: key}
}

/*
Get the results channel.

You must listen on the results channel (even if only to throw them away)
or otherwise all the workers will deadlock.
*/
func (b *BulkWriter) Results() <-chan Result {
	return b.results
}

/*
Close this BulkWriter, and wait until all our existing operations have
completed. You must not call Write anymore after CloseWait has been called.

CloseWait will panic if called more than once on a BulkWriter.
*/
func (b *BulkWriter) CloseWait() {
	close(b.ch)
	b.wg.Wait()
	close(b.results)
}

func (b *BulkWriter) main(perWrite int) {
	defer func() {
		b.wg.Done()
	}()
	var running = true
	for running {
		g := group{}

		for i := 0; i < perWrite; i++ {
			var msg message
			msg, running = <-b.ch
			if running {
				if msg.doc != nil {
					g.docs = append(g.docs, msg.doc)
				} else {
					g.deleteKeys = append(g.deleteKeys, msg.deleteKey)
				}
			} else {
				break
			}
		}
		if g.length() > 0 {
			b.groups <- g
		}
	}
	close(b.groups)
}

func (b *BulkWriter) worker(id int) {
	defer func() {
		b.wg.Done()
	}()
	for group := range b.groups {
		origGroup := group
		waitFor := 100 * time.Millisecond
		for i := 0; i < 5; i++ {
			group = b.runBatch(group, &waitFor)
			if group.length() < origGroup.length() {
				break
			}
		}

		// Run any remaining documents individually
		for _, doc := range group.docs {
			rlist := []dynago.Document{doc}
			for {
				_, err := b.client.PutItem(b.table, doc).Execute()
				if err == nil {
					b.results <- Result{Documents: rlist}
					break
				} else if !b.retryLogic(err, &waitFor, rlist, nil) {
					break
				}
			}
		}
		// run any remaining deletes individually
		for _, key := range group.deleteKeys {
			dList := []dynago.Document{key}
			_, err := b.client.DeleteItem(b.table, key).Execute()
			if err == nil {
				b.results <- Result{DeleteKeys: dList}
				break
			} else if !b.retryLogic(err, &waitFor, nil, dList) {
				break
			}
		}
	}
}

func (b *BulkWriter) runBatch(g group, waitFor *time.Duration) group {
	batch := b.client.BatchWrite()
	if len(g.docs) > 0 {
		batch = batch.Put(b.table, g.docs...)
	}
	if len(g.deleteKeys) > 0 {
		batch = batch.Delete(b.table, g.deleteKeys...)
	}
	result, err := batch.Execute()
	if err == nil {
		b.results <- Result{Documents: g.docs, DeleteKeys: g.deleteKeys}
		g = group{}
		for _, item := range result.UnprocessedItems[b.table] {
			if item.PutRequest != nil {
				g.docs = append(g.docs, item.PutRequest.Item)
			} else {
				g.deleteKeys = append(g.deleteKeys, item.DeleteRequest.Key)
			}
		}
	} else if !b.retryLogic(err, waitFor, g.docs, g.deleteKeys) {
		return group{}
	}
	return g
}

func (b *BulkWriter) retryLogic(err error, waitFor *time.Duration, toWrite []dynago.Document, toDelete []dynago.Document) bool {
	if e, ok := err.(*dynago.Error); ok {
		if canRetry(e) {
			time.Sleep(*waitFor)
			*waitFor *= 2
			return true
		} else {
			b.results <- Result{Documents: toWrite, DeleteKeys: toDelete, Error: err, DynagoError: e}
			return false
		}
	} else {
		b.results <- Result{Documents: toWrite, DeleteKeys: toDelete, Error: err}
		return false
	}
}

func canRetry(e *dynago.Error) bool {
	switch e.Type {
	case dynago.ErrorThrottling, dynago.ErrorThroughputExceeded:
		return true
	case dynago.ErrorServiceUnavailable, dynago.ErrorInternalFailure:
		return true
	default:
		return false
	}
}

type group struct {
	docs       []dynago.Document
	deleteKeys []dynago.Document
}

func (g group) length() int {
	return len(g.docs) + len(g.deleteKeys)
}

type message struct {
	doc       dynago.Document
	deleteKey dynago.Document
}

type Result struct {
	Documents   []dynago.Document // The documents we're talking about
	DeleteKeys  []dynago.Document // Deleted keys
	Error       error             // If there's an error, then this is set
	DynagoError *dynago.Error     // If the error happens to be a dynago.Error, then we set this too.
}
