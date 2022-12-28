package zfs

import (
	"context"
	"fmt"
)

type DatasetFilter interface {
	Filter(p *DatasetPath) (pass bool, err error)
	ConfigErrors(ctx context.Context) (count float64)
}

// Returns a DatasetFilter that does not filter (passes all paths)
func NoFilter() DatasetFilter {
	return noFilter{}
}

type noFilter struct{}

var _ DatasetFilter = noFilter{}

func (noFilter) Filter(p *DatasetPath) (pass bool, err error) { return true, nil }
func (noFilter) ConfigErrors(ctx context.Context) (count float64) { return 0 }

func ZFSListMapping(ctx context.Context, filter DatasetFilter) (datasets []*DatasetPath, err error) {
	res, err := ZFSListMappingProperties(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	datasets = make([]*DatasetPath, len(res))
	for i, r := range res {
		datasets[i] = r.Path
	}
	return datasets, nil
}

type ZFSListMappingPropertiesResult struct {
	Path *DatasetPath
	// Guaranteed to have the same length as properties in the originating call
	Fields []string
}

// properties must not contain 'name'
func ZFSListMappingProperties(ctx context.Context, filter DatasetFilter, properties []string) (datasets []ZFSListMappingPropertiesResult, err error) {

	if filter == nil {
		panic("filter must not be nil")
	}

	for _, p := range properties {
		if p == "name" {
			panic("properties must not contain 'name'")
		}
	}
	newProps := make([]string, len(properties)+1)
	newProps[0] = "name"
	copy(newProps[1:], properties)
	properties = newProps

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	rchan := make(chan ZFSListResult)

	go ZFSListChan(ctx, rchan, properties, nil, "-r", "-t", "filesystem,volume")

	datasets = make([]ZFSListMappingPropertiesResult, 0)
	for r := range rchan {

		if r.Err != nil {
			err = r.Err
			return
		}

		var path *DatasetPath
		if path, err = NewDatasetPath(r.Fields[0]); err != nil {
			return
		}

		pass, filterErr := filter.Filter(path)
		if filterErr != nil {
			return nil, fmt.Errorf("error calling filter: %s", filterErr)
		}
		if pass {
			datasets = append(datasets, ZFSListMappingPropertiesResult{
				Path:   path,
				Fields: r.Fields[1:],
			})
		}

	}

	return
}
