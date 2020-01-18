package controllers

import (
	"fmt"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func hasValue(s interface{}, elem interface{}) bool {
	arrV := reflect.ValueOf(s)

	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {
			if arrV.Index(i).Interface() == elem {
				return true
			}
		}
	}
	return false
}

func NamespaceNameKey(meta *metav1.ObjectMeta) string {
	return fmt.Sprintf("%s/%s", meta.Namespace, meta.Name)
}
