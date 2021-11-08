package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
)

func WebhookCallback(c *gin.Context) {
	var body []byte
	if c.Request.Body != nil {
		if data, err := ioutil.ReadAll(c.Request.Body); err == nil {
			body = data
		}
	}

	if len(body) == 0 {
		glog.Error("empty body")
		return
	}

	var admissionResponse *v1beta1.AdmissionResponse

	ar := v1beta1.AdmissionReview{}

	admissionReview := v1beta1.AdmissionReview{}

	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		glog.Errorf("Can't decode body: %v", err)
		admissionResponse = &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	admissionResponse = mutate(&ar)
	admissionReview.Response = admissionResponse

	resp, _ := json.Marshal(admissionReview)

	c.Writer.Write(resp)
}

func mutate(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {

	return &v1beta1.AdmissionResponse{
		Allowed: true,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}
