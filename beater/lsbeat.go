package beater

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/beats/v7/packetbeat/pb"

	"github.com/JustinChen1975/lsbeat/config"
)

// lsbeat configuration.
type lsbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
	lastIndexTime time.Time
}

// New creates an instance of lsbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	//c := DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &lsbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts lsbeat.
func (bt *lsbeat) Run(b *beat.Beat) error {
	logp.Info("lsbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	for {
		now := time.Now()
		//bt.listDir(bt.config.Path, b.Info.Beat) // call listDir
		bt.listDir1(bt.config.Path, b.Info.Beat) // call listDir
		bt.lastIndexTime = now                  // mark Timestamp
		logp.Info("Event sent")
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}
	}
	return nil

	//ticker := time.NewTicker(bt.config.Period)
	//counter := 1
	//for {
	//	select {
	//	case <-bt.done:
	//		return nil
	//	case <-ticker.C:
	//	}
	//
	//	event := beat.Event{
	//		Timestamp: time.Now(),
	//		Fields: common.MapStr{
	//			"type":    b.Info.Name,
	//			"counter": counter,
	//		},
	//	}
	//	bt.client.Publish(event)
	//	logp.Info("Event sent")
	//	counter++
	//}
}

// Stop stops lsbeat.
func (bt *lsbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

//only for testing

func (bt *lsbeat) listDir1(dirFile string, beatname string) {
	files, _ := ioutil.ReadDir(dirFile)
	for _, f := range files {
		t := f.ModTime()
		path := filepath.Join(dirFile, f.Name())
		if t.After(bt.lastIndexTime) {

			evt, pbf := pb.NewBeatEvent(t)

			//pbf.SetSource(&t.src)
			//pbf.SetDestination(&t.dst)
			//pbf.Network.Transport = t.transport.String()
			pbf.Network.Protocol = "dns"
			//pbf.Error.Message = t.notes

			fields := evt.Fields
			fields["type"] = beatname
			fields["modtime"] =common.Time(t)

			//fields["filename"] =f.Name()
			//fields["path"] =path
			//fields["directory"] =f.IsDir()
			//fields["filesize"] =f.Size()

			lsEvent := common.MapStr{}
			fields["listDirectory"] = lsEvent

			lsEvent["filename"] =f.Name()
			lsEvent["path"] =path
			lsEvent["directory"] =f.IsDir()
			lsEvent["filesize"] =f.Size()

			//fields["status"] = common.

			//event := beat.Event{
			//	Timestamp: time.Now(),
			//	Fields: common.MapStr {
			//		"type":       beatname,
			//		"modtime":    common.Time(t),
			//		"filename":   f.Name(),
			//		"path":       path,
			//		"directory":  f.IsDir(),
			//		"filesize":   f.Size(),
			//	},
			//}

			bt.client.Publish(evt)
		}
		if f.IsDir() {
			bt.listDir(path, beatname)
		}
	}
}

func (bt *lsbeat) listDir(dirFile string, beatname string) {
	files, _ := ioutil.ReadDir(dirFile)
	for _, f := range files {
		t := f.ModTime()
		path := filepath.Join(dirFile, f.Name())
		if t.After(bt.lastIndexTime) {

			event := beat.Event{
				Timestamp: time.Now(),
				Fields: common.MapStr {
					"type":       beatname,
					"modtime":    common.Time(t),
					"filename":   f.Name(),
					"path":       path,
					"directory":  f.IsDir(),
					"filesize":   f.Size(),
				},
			}

			bt.client.Publish(event)
		}
		if f.IsDir() {
			bt.listDir(path, beatname)
		}
	}
}
