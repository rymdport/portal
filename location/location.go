// Package location lets sandboxed applications query basic information about the location.
package location

const interfaceName = "org.freedesktop.portal.Location"

// Location contains the returned location data.
type Location struct {
	Latitude  float64   // The latitude, in degrees.
	Longitude float64   // The longitude, in degrees.
	Altitude  float64   // The altitude, in meters.
	Accuracy  float64   // The accuracy, in meters.
	Speed     float64   // The speed, in meters per second.
	Heading   float64   // The heading, in degrees, going clockwise. North 0, East 90, South 180, West 270.
	Timestamp [2]uint64 // The timestamp, as seconds and microseconds since the Unix epoch.
}
