// Package pipeline provides a composable stage-based processing pipeline
// for port activity diffs.
//
// WithSummarizer stage accumulates opened/closed port counts over a scan
// interval and flushes a periodic summary via the provided Summarizer.
// It is typically placed near the end of the pipeline so all filtering
// and enrichment stages have already run.
//
// Example usage:
//
//	sum := summarizer.New(30*time.Second, func(s summarizer.Summary) {
//		fmt.Printf("opened=%d closed=%d\n", s.Opened, s.Closed)
//	})
//	p := pipeline.New(
//		pipeline.WithFilter(f),
//		pipeline.WithSummarizer(sum),
//	)
package pipeline
