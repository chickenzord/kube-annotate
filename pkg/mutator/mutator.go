package mutator

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/chickenzord/kube-annotate/pkg/config"
	"k8s.io/api/admission/v1beta1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

//Patch patching operation
type Patch struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()

	// (https://github.com/kubernetes/kubernetes/issues/57982)
	defaulter = runtime.ObjectDefaulter(runtimeScheme)
)

func init() {
	_ = corev1.AddToScheme(runtimeScheme)
	_ = admissionregistrationv1beta1.AddToScheme(runtimeScheme)
	// defaulting with webhooks:
	// https://github.com/kubernetes/kubernetes/issues/57982
	_ = v1.AddToScheme(runtimeScheme)
}

func respond(review *v1beta1.AdmissionReview, response *v1beta1.AdmissionResponse) *v1beta1.AdmissionReview {
	result := &v1beta1.AdmissionReview{}
	if response != nil {
		result.Response = response
		if review.Request != nil {
			result.Response.UID = review.Request.UID
		}
	}
	return result
}

func respondWithError(review *v1beta1.AdmissionReview, err error) *v1beta1.AdmissionReview {
	return respond(review, &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	})
}

func respondWithSkip(review *v1beta1.AdmissionReview) *v1beta1.AdmissionReview {
	return respond(review, &v1beta1.AdmissionResponse{
		Allowed: true,
	})
}

func respondWithPatches(review *v1beta1.AdmissionReview, patches []Patch) *v1beta1.AdmissionReview {
	patchesBytes, err := json.Marshal(patches)
	if err != nil {
		return respondWithError(review, fmt.Errorf("cannot serialize patches: %v", err))
	}

	return respond(review, &v1beta1.AdmissionResponse{
		Allowed: true,
		Patch:   patchesBytes,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	})
}

func createPatchFromAnnotations(base, extra map[string]string) Patch {
	if base == nil {
		return Patch{
			Op:    "add",
			Path:  "/metadata/annotations",
			Value: extra,
		}
	}

	annotations := make(map[string]string)
	for k, v := range base {
		annotations[k] = v
	}
	if extra != nil {
		for k, v := range extra {
			annotations[k] = v
		}
	}

	return Patch{
		Op:    "replace",
		Path:  "/metadata/annotations",
		Value: annotations,
	}
}

//ParseBytes parses AdmissionReview into bytes
func ParseBytes(data []byte) (*v1beta1.AdmissionReview, error) {
	result := v1beta1.AdmissionReview{}
	if _, _, err := deserializer.Decode(data, nil, &result); err != nil {
		return nil, fmt.Errorf("cannot deserialize data to AdmissionReview: %v", err)
	}

	return &result, nil
}

//MutateBytes mutates AdmissionReview bytes
func MutateBytes(data []byte) ([]byte, error) {
	reviewRequest := v1beta1.AdmissionReview{}
	if _, _, err := deserializer.Decode(data, nil, &reviewRequest); err != nil {
		return nil, fmt.Errorf("cannot decode data to AdmissionReview: %v", err)
	}

	reviewResponse := Mutate(&reviewRequest)

	result, err := json.Marshal(reviewResponse)
	if err != nil {
		return nil, fmt.Errorf("cannot encode response: %v", err)
	}

	return result, nil
}

//Mutate mutates AdmissionReview
func Mutate(review *v1beta1.AdmissionReview) *v1beta1.AdmissionReview {
	//deserialize pod
	var pod corev1.Pod
	if err := json.Unmarshal(review.Request.Object.Raw, &pod); err != nil {
		log.WithData(review).WithError(err).Errorf("error mutating pod")
		return respondWithError(review, errors.New("cannot deserialize pod from AdmissionRequest"))
	}

	//create patches based on rules
	log.WithData(review).Debug("processing AdmissionReview")
	patches := make([]Patch, 0)
	for _, rule := range config.Rules {
		if rule.Selector.AsSelector().Matches(labels.Set(pod.Labels)) {
			patch := createPatchFromAnnotations(pod.Annotations, rule.Annotations)
			patches = append(patches, patch)
		}
	}

	if len(patches) > 0 {
		log.WithData(review).Infof("mutating Pod with %d patch(es)", len(patches))
		return respondWithPatches(review, patches)
	}

	log.Infof("skipping Pod")
	return respondWithSkip(review)
}
