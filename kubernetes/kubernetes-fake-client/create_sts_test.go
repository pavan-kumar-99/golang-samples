package sts

import (
	"context"
	"fmt"
	"testing"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestCreate_sts(t *testing.T) {
	type cases []struct {
		namespace string
		name      string
		replicas  int32
	}

	values := cases{
		{
			name:      "go-sts",
			namespace: "default",
			replicas:  2,
		},
		{
			name:      "go-sts1",
			namespace: "default1",
			replicas:  2,
		},
		{
			name:      "go-sts2",
			namespace: "default1",
			replicas:  2,
		},
		{
			name:      "go-sts3",
			namespace: "default1",
			replicas:  2,
		},
		{
			name:      "go-sts4",
			namespace: "default1",
			replicas:  2,
		},
		{
			name:      "go-sts5",
			namespace: "default1",
			replicas:  2,
		},
		{
			name:      "go-sts6",
			namespace: "default1",
			replicas:  2,
		},
	}

	api := &StsSpec{
		Client: testclient.NewSimpleClientset(),
	}

	for _, c := range values {
		testname := fmt.Sprintf("%s,%s", c.name, c.namespace)
		t.Run(testname, func(t *testing.T) {
			err := api.Create_sts(c.namespace, c.name, c.replicas)
			if err != nil {
				t.Fatal(err.Error())
			}

			_, err = api.Client.AppsV1().StatefulSets(c.namespace).Get(context.TODO(), c.name, v1.GetOptions{})
			if err != nil {
				t.Fatal(err.Error())
			}
		})
	}
}
