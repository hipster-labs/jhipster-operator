package pkg

import (
	"bytes"
	"html/template"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetK8sObjectFromTemplate(fileName string, vars any) (runtime.Object, *schema.GroupVersionKind, error) {
	temp, err := template.ParseFiles(fileName)

	if err != nil {
		return nil, nil, err
	}
	var data bytes.Buffer

	err = temp.Execute(&data, vars)

	if err != nil {
		return nil, nil, err
	}

	return GetK8sObjectFromBytes(data.Bytes())
}

func GetK8sObjectFromBytes(data []byte) (runtime.Object, *schema.GroupVersionKind, error) {
	decoder := scheme.Codecs.UniversalDeserializer()
	obj, groupVersionKind, err := decoder.Decode(data, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	return obj, groupVersionKind, nil
}
