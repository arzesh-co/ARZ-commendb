package CommenDb

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type filter struct {
	Condition any
	Label     string
	Operation string
}
type aggregation struct {
	GroupBy     string `json:"group_by"`
	Aggregators []struct {
		Aggregate string `json:"aggregate"`
		Operation string `json:"operation"`
	} `json:"aggregators"`
}
type sort struct {
	DbName string `json:"db_name"`
	Type   string `json:"type"`
}

func createFilter(cond filter) interface{} {
	switch cond.Operation {
	case "Start With":
		return primitive.Regex{Pattern: "^" + cond.Condition.(string) + ".", Options: "i"}
	case "End With":
		return primitive.Regex{Pattern: ".*" + cond.Condition.(string) + "$", Options: "i"}
	case "Equal":
		return bson.M{"$eq": cond.Condition}
	case "Include":
		return primitive.Regex{Pattern: ".*" + cond.Condition.(string) + ".*", Options: "i"}
	case "Empty":
		return bson.M{"$exists": false}
	case "not Empty":
		return bson.M{"$exists": true}
	case "=":
		return bson.M{"$eq": ConvertFilterCondition(cond.Condition)}
	case ">=":
		return bson.M{"$gte": ConvertFilterCondition(cond.Condition)}
	case "<=":
		return bson.M{"$lte": ConvertFilterCondition(cond.Condition)}
	case ">":
		return bson.M{"$gt": ConvertFilterCondition(cond.Condition)}
	case "<":
		return bson.M{"$lt": ConvertFilterCondition(cond.Condition)}
	}
	return bson.M{}
}
func ConvertFilterCondition(condition any) any {
	switch condition.(type) {
	case string:
		switch ConvertorType(condition.(string)) {
		case "func":
			return findFunc(condition.(string))
		case "string":
			return condition
		default:
			return condition
		}
	default:
		return condition
	}
}
func CreateAggregation(aggr string) map[string]interface{} {
	agg := &aggregation{}
	err := json.Unmarshal([]byte(aggr), agg)
	if err != nil {
		return nil
	}
	filter := make(map[string]interface{})
	if aggr == "" {
		return nil
	}
	filter["_id"] = "$" + agg.GroupBy
	for _, aggregator := range agg.Aggregators {
		switch aggregator.Operation {
		case "avg":
			filter[aggregator.Aggregate] = bson.M{"$avg": "$" + aggregator.Aggregate}
		case "sum":
			filter[aggregator.Aggregate] = bson.M{"$sum": "$" + aggregator.Aggregate}
		case "count":
			filter[aggregator.Aggregate] = bson.M{"$sum": 1}
		case "min":
			filter[aggregator.Aggregate] = bson.M{"$min": "$" + aggregator.Aggregate}
		case "max":
			filter[aggregator.Aggregate] = bson.M{"$max": "$" + aggregator.Aggregate}
		}
	}
	return filter
}
func CreateFilter(filterString string) (map[string]any, error) {
	var filters []filter
	err := json.Unmarshal([]byte(filterString), &filters)
	if err != nil {
		return nil, err
	}
	clintFilterMap := make(map[string]any)
	for _, f := range filters {
		clintFilterMap[f.Label] = createFilter(f)
	}
	return clintFilterMap, nil
}
func CreateSorting(sortString string) (map[string]any, error) {
	var sorts []sort
	err := json.Unmarshal([]byte(sortString), &sorts)
	if err != nil {
		return nil, err
	}
	sortFilter := make(map[string]any)
	for _, s := range sorts {
		switch s.Type {
		case "asc":
			sortFilter[s.DbName] = 1
		case "des":
			sortFilter[s.DbName] = -1
		}
	}
	return sortFilter, nil
}
