// Package pipeline provides a composable stage-based processing pipeline
// for port activity diffs.
//
// # GeoIP Stage
//
// WithGeoIP enriches each opened port entry with geographic metadata by
// performing a remote IP geolocation lookup. Private and loopback addresses
// are skipped automatically.
//
// Results are cached inside the [geoip.Client] to avoid redundant HTTP
// requests across pipeline runs.
//
// Usage:
//
//	geo, err := geoip.New(geoip.DefaultConfig())
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	p := pipeline.New(
//		pipeline.WithGeoIP(geo),
//	)
package pipeline
