package lib

import (
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strings"
)

func isZeroOfUnderlyingType(i interface{}) bool {
	return reflect.DeepEqual(i, reflect.Zero(reflect.TypeOf(i)).Interface())
}

func PrettyBsonSet(field string, model interface{}) bson.M {
	var update bson.M = map[string]interface{}{}

	val := reflect.ValueOf(model).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		v := valueField.Interface()

		typeField := val.Type().Field(i)
		tag := typeField.Tag
		bsonTags := strings.Split(tag.Get("bson"), ",")

		if Contains(bsonTags, "omitempty") && isZeroOfUnderlyingType(v) {
			continue
		}

		bsonField := bsonTags[0]

		update[field+"."+bsonField] = v
	}

	return update
}
