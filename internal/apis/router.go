package apis

import (
	"slices"
	"sync"

	"github.com/godbus/dbus/v5"
)

// Signal routing. godbus delivers every signal on the connection to every
// channel registered via conn.Signal, which is wasteful for portal workloads
// and prone to leaks when callers forget conn.RemoveSignal. We register one
// channel with godbus and fan out here, keyed by (path, interface.member).

type routerKey struct {
	path dbus.ObjectPath
	name string // interface.member
}

type router struct {
	attach sync.Once

	mu   sync.Mutex
	subs map[routerKey][]chan<- *dbus.Signal
}

var defaultRouter = &router{subs: map[routerKey][]chan<- *dbus.Signal{}}

func (r *router) attachTo(conn *dbus.Conn) {
	r.attach.Do(func() {
		in := make(chan *dbus.Signal, 256)
		conn.Signal(in)
		go r.loop(in)
	})
}

func (r *router) loop(in <-chan *dbus.Signal) {
	for sig := range in {
		key := routerKey{path: sig.Path, name: sig.Name}

		// snapshot the subscriber list so we don't hold the lock while sending
		r.mu.Lock()
		targets := append([]chan<- *dbus.Signal(nil), r.subs[key]...)
		r.mu.Unlock()

		for _, ch := range targets {
			select {
			case ch <- sig:
			default: // slow subscriber, drop
			}
		}
	}
}

func (r *router) subscribe(key routerKey, ch chan<- *dbus.Signal) (cleanup func()) {
	r.mu.Lock()
	r.subs[key] = append(r.subs[key], ch)
	r.mu.Unlock()

	return func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		subs := r.subs[key]
		if i := slices.Index(subs, ch); i >= 0 {
			subs = slices.Delete(subs, i, i+1)
		}
		if len(subs) == 0 {
			delete(r.subs, key)
		} else {
			r.subs[key] = subs
		}
	}
}

// SubscribeSignal returns a channel that receives signals matching
// (path, interface.member), and a cleanup to unsubscribe. The caller is
// still responsible for AddMatchSignal/RemoveMatchSignal so the bus
// forwards the signals to this connection.
//
// conn must be the shared session bus: the router attaches to whichever
// connection it sees first and ignores any other.
func SubscribeSignal(conn *dbus.Conn, path dbus.ObjectPath, interfaceName, memberName string) (<-chan *dbus.Signal, func()) {
	defaultRouter.attachTo(conn)

	key := routerKey{path: path, name: interfaceName + "." + memberName}
	ch := make(chan *dbus.Signal, 16) // absorbs bursts like SettingChanged on theme switch
	return ch, defaultRouter.subscribe(key, ch)
}
