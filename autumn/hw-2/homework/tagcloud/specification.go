package tagcloud

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	stats        []*TagStat
	tags         map[string]int
	maxOccurence int
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
// TODO: You decide whether this function should return a pointer or a value
func New() *TagCloud {
	return &TagCloud{
		stats:        make([]*TagStat, 0),
		tags:         make(map[string]int),
		maxOccurence: 1,
	}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
// TODO: You decide whether receiver should be a pointer or a value
func (t *TagCloud) AddTag(tag string) {
	idx, ok := t.tags[tag]
	if !ok {
		idx := len(t.stats)
		t.stats = append(t.stats, &TagStat{
			Tag:             tag,
			OccurrenceCount: 1,
		})
		t.tags[tag] = idx
		return
	}
	t.stats[idx].OccurrenceCount++
	if t.maxOccurence < t.stats[idx].OccurrenceCount {
		t.maxOccurence = t.stats[idx].OccurrenceCount
	}
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
// TODO: You decide whether receiver should be a pointer or a value
func (t *TagCloud) TopN(n int) []*TagStat {
	if n > len(t.stats) {
		return t.stats
	}
	freq := make([][]int, t.maxOccurence)
	for _, i := range t.tags {
		idx := t.stats[i].OccurrenceCount - 1
		freq[idx] = append(freq[idx], i)
	}
	res := make([]*TagStat, 0)
	for i := len(freq) - 1; i >= 0; i-- {
		for _, idx := range freq[i] {
			res = append(res, t.stats[idx])
			if len(res) == n {
				return res
			}
		}
	}
	return nil
}
