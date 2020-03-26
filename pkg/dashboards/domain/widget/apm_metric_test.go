package widget

import "testing"

type testCase struct {
	first  ApmMetric
	second ApmMetric
}

func equalTestCases() []testCase {
	return []testCase{
		{
			first: ApmMetric{
				Duration:  0,
				EntityIds: []int{},
				Metrics:   nil,
				Facet:     "",
				OrderBy:   "",
			},
			second: ApmMetric{
				Duration:  0,
				EntityIds: []int{},
				Metrics:   nil,
				Facet:     "",
				OrderBy:   "",
			},
		},
		{
			first: ApmMetric{
				Duration:  0,
				EntityIds: []int{2, 1},
				Metrics:   nil,
				Facet:     "facet",
				OrderBy:   "order",
			},
			second: ApmMetric{
				Duration:  0,
				EntityIds: []int{1, 2},
				Metrics:   nil,
				Facet:     "facet",
				OrderBy:   "order",
			},
		},
		{
			first: ApmMetric{
				Duration:  0,
				EntityIds: []int{},
				Metrics: MetricList{
					{
						Name:   "metric 1",
						Values: []string{"value 1"},
					},
					{
						Name:   "metric 2",
						Values: []string{"value 2"},
					},
				},
				Facet:   "",
				OrderBy: "",
			},
			second: ApmMetric{
				Duration:  0,
				EntityIds: []int{},
				Metrics: MetricList{
					{
						Name:   "metric 2",
						Values: []string{"value 2"},
					},
					{
						Name:   "metric 1",
						Values: []string{"value 1"},
					},
				},
				Facet:   "",
				OrderBy: "",
			},
		},
	}
}

func TestApmMetric_Equals(t *testing.T) {
	for i, testCase := range equalTestCases() {
		if !testCase.first.Equals(&testCase.second) {
			t.Errorf("Metrics in case %d must be equal", i+1)
		}
	}
}

func notEqualTestCases() []testCase {
	return []testCase{
		{
			first: ApmMetric{
				Duration:  10,
				EntityIds: []int{},
				Metrics:   nil,
				Facet:     "",
				OrderBy:   "",
			},
			second: ApmMetric{
				Duration:  0,
				EntityIds: []int{},
				Metrics:   nil,
				Facet:     "",
				OrderBy:   "",
			},
		},
		{
			first: ApmMetric{
				Duration:  1,
				EntityIds: []int{3},
				Metrics:   nil,
				Facet:     "",
				OrderBy:   "",
			},
			second: ApmMetric{
				Duration:  0,
				EntityIds: []int{},
				Metrics:   nil,
				Facet:     "",
				OrderBy:   "",
			},
		},
		{
			first: ApmMetric{
				Duration:  0,
				EntityIds: []int{},
				Metrics:   nil,
				Facet:     "facet",
				OrderBy:   "",
			},
			second: ApmMetric{
				Duration:  0,
				EntityIds: []int{},
				Metrics:   nil,
				Facet:     "",
				OrderBy:   "",
			},
		},
		{
			first: ApmMetric{
				Duration:  0,
				EntityIds: []int{},
				Metrics:   MetricList{
					{
						Name:   "metric 1",
						Values: nil,
					},
				},
				Facet:     "",
				OrderBy:   "",
			},
			second: ApmMetric{
				Duration:  0,
				EntityIds: []int{},
				Metrics:   MetricList{
					{
						Name:   "metric 2",
						Values: nil,
					},
				},
				Facet:     "",
				OrderBy:   "",
			},
		},
		{
			first: ApmMetric{
				Duration:  0,
				EntityIds: []int{},
				Metrics:   MetricList{
					{
						Name:   "metric 1",
						Values: []string{"value 1"},
					},
				},
				Facet:     "",
				OrderBy:   "",
			},
			second: ApmMetric{
				Duration:  0,
				EntityIds: []int{},
				Metrics:   MetricList{
					{
						Name:   "metric 1",
						Values: []string{"value 2"},
					},
				},
				Facet:     "",
				OrderBy:   "",
			},
		},
	}
}

func TestApmMetric_NotEquals(t *testing.T) {
	for i, testCase := range notEqualTestCases() {
		if testCase.first.Equals(&testCase.second) {
			t.Errorf("Metrics in case %d should not be equal", i+1)
		}
	}
}
