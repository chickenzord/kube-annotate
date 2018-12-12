package annotator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/chickenzord/kube-annotate/config"
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

func mutationRequired(metadata *metav1.ObjectMeta) bool {
	return false
}

func parseBody(r *http.Request) (*v1beta1.AdmissionReview, error) {
	if r.ContentLength == 0 {
		return nil, errors.New("Empty Body")
	}

	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		return nil, fmt.Errorf("Invalid content type: %s", contentType)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("Cannot read body: %v", err)
	}

	result := v1beta1.AdmissionReview{}
	if _, _, err := deserializer.Decode(data, nil, &result); err != nil {
		return nil, fmt.Errorf("Cannot deserialize: %v", err)
	}

	return &result, nil
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
	log.Errorf("Error mutating Pod: %v", err)
	return respond(review, &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	})
}

func respondWithSkip(review *v1beta1.AdmissionReview) *v1beta1.AdmissionReview {
	log.Infof("Skipping Pod")
	return respond(review, &v1beta1.AdmissionResponse{
		Allowed: true,
	})
}

func respondWithPatches(review *v1beta1.AdmissionReview, patches []Patch) *v1beta1.AdmissionReview {
	patchesBytes, err := json.Marshal(patches)
	if err != nil {
		return respondWithError(review, errors.New("Cannot serialize patches to JSON"))
	}

	log.Infof("Mutating Pod with %d patch(es)", len(patches))
	return respond(review, &v1beta1.AdmissionResponse{
		Allowed: true,
		Patch:   patchesBytes,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	})
}

func createPatchesFromAnnotations(base, extra map[string]string) []Patch {
	patches := make([]Patch, 0)
	for key, val := range extra {
		if base == nil || base[key] == "" {
			base = map[string]string{}
			patches = append(patches, Patch{
				Op:   "add",
				Path: "/metadata/annotations",
				Value: map[string]string{
					key: val,
				},
			})
		} else {
			patches = append(patches, Patch{
				Op:    "replace",
				Path:  "/metadata/annotations/" + key,
				Value: val,
			})
		}
	}
	return patches
}

func mutate(review *v1beta1.AdmissionReview) *v1beta1.AdmissionReview {
	var pod corev1.Pod
	if err := json.Unmarshal(review.Request.Object.Raw, &pod); err != nil {
		return respondWithError(review, errors.New("Cannot deserialize Pod from AdmissionRequest"))
	}

	// log
	podObject, err := json.Marshal(pod)
	if err != nil {
		log.Errorf("Cannot serialize Pod: %v", err)
	}
	log.Debug(string(podObject))

	for _, rule := range config.Rules {
		if rule.Selector.AsSelector().Matches(labels.Set(pod.Labels)) {
			return respondWithPatches(review, createPatchesFromAnnotations(pod.Annotations, rule.Annotations))
		}
	}

	return respondWithSkip(review)
}
