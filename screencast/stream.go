package screencast

import (
	"github.com/godbus/dbus/v5"
)

// PipeWire stream
type Stream struct {
	NodeId     uint32
	Properties map[string]dbus.Variant
}

func (s *Stream) Id() (string, bool) {
	val, ok := s.Properties["id"]
	if !ok {
		return "", ok
	}
	return val.Value().(string), true
}

func (s *Stream) Position() (int32, int32, bool) {
	val, ok := s.Properties["position"]
	if !ok {
		return 0, 0, false
	}

	position, ok := val.Value().([]interface{})
	if len(position) != 2 || !ok {
		return 0, 0, false
	}

	return position[0].(int32), position[1].(int32), true
}

func (s *Stream) Size() (int32, int32, bool) {
	val, ok := s.Properties["size"]
	if !ok {
		return 0, 0, false
	}

	size, ok := val.Value().([]interface{})
	if len(size) != 2 || !ok {
		return 0, 0, false
	}

	return size[0].(int32), size[1].(int32), true
}

func (s *Stream) SourceType() (SourceType, bool) {
	val, ok := s.Properties["source_type"]
	if !ok {
		return 0, false
	}

	return SourceType(val.Value().(uint32)), true
}

func (s *Stream) MappingId() (string, bool) {
	val, ok := s.Properties["mapping_id"]
	if !ok {
		return "", false
	}

	return val.Value().(string), true
}

func (s *Stream) PipeWireSerial() (uint64, bool) {
	val, ok := s.Properties["pipewire-serial"]
	if !ok {
		return 0, false
	}

	return val.Value().(uint64), true
}
