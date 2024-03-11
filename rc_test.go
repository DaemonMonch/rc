package rc

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

type tt struct{}

func (tt) Unmarshall() (interface{}, error) {
	b, err := os.ReadFile("./test.cfg")
	if len(b) == 0 {
		return nil, io.EOF
	}
	return string(b), err
}

func TestRc(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rc := NewReloadableConfig(ctx, NewFsNotifyPoller("./test.cfg", fsnotify.Write), &tt{})
	i := 0
	for {
		c, ok := <-rc.Notify
		if i > 10 {
			cancel()
		}

		if !ok {
			fmt.Println("ffff")
			cancel()
			return
		}
		fmt.Println(c.(string))
		i++
	}
}

type TestYamlConfig struct {
	A string
	B struct {
		Ba  string
		Bas []string
	}
}

func TestRcYaml(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rc := NewReloadableConfig(ctx, NewFsNotifyPoller("./test.yml", fsnotify.Write), NewYamlConfig("./test.yml", new(TestYamlConfig)))
	i := 0
	for {
		c, ok := <-rc.Notify
		if i > 10 {
			cancel()
		}

		if !ok {
			fmt.Println("ffff")
			cancel()
			return
		}
		fmt.Printf("%v \n", c.(*TestYamlConfig))
		i++
	}
}

func TestRcTimedYaml(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	timer := time.NewTicker(1 * time.Second)
	rc := NewReloadableConfig(ctx, (*TimedPoller)(timer), NewYamlConfig("./test.yml", new(TestYamlConfig)))
	i := 0
	for {
		c, ok := <-rc.Notify
		if i > 10 {
			cancel()
		}

		if !ok {
			fmt.Println("ffff")
			cancel()
			return
		}
		fmt.Printf("%v \n", c.(*TestYamlConfig))
		i++
	}
}
