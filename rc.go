package rc

import (
	"context"
	"log"
)

type ConfigUnmarshaller interface {
	Unmarshall() (interface{}, error)
}

type ReloadableConfig struct {
	Poller             Poller
	ConfigUnMarshaller ConfigUnmarshaller
	Notify             chan interface{}
}

func NewReloadableConfig(ctx context.Context, poller Poller, conConfigUnMarshaller ConfigUnmarshaller) *ReloadableConfig {
	rc := &ReloadableConfig{poller, conConfigUnMarshaller, make(chan interface{}, 1)}
	err := rc.marshallConfigAndNotify()
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		defer rc.Poller.Close()
		defer close(rc.Notify)
		for {
			select {
			case evt := <-rc.Poller.Poll():
				_ = evt
				rc.marshallConfigAndNotify()
			case <-ctx.Done():
				return
			}
		}
	}()

	return rc
}

func (rc *ReloadableConfig) marshallConfigAndNotify() error {
	newConfig, err := rc.ConfigUnMarshaller.Unmarshall()
	if err != nil {
		log.Printf("Unmarshall config fail , err = %v \n", err)
		return err
	}
	select {
	case rc.Notify <- newConfig:
		break
	default:
		log.Println("notify chan busy, ignore")
	}

	return nil
}
