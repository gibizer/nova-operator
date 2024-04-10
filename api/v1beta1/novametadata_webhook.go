/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

//
// Generated by:
//
// operator-sdk create webhook --group nova --version v1beta1 --kind NovaMetadata --programmatic-validation --defaulting
//

package v1beta1

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// NovaMetadataDefaults -
type NovaMetadataDefaults struct {
	ContainerImageURL string
}

var novaMetadataDefaults NovaMetadataDefaults

// log is for logging in this package.
var novametadatalog = logf.Log.WithName("novametadata-resource")

// SetupNovaMetadataDefaults - initialize NovaMetadata spec defaults for use with either internal or external webhooks
func SetupNovaMetadataDefaults(defaults NovaMetadataDefaults) {
	novaMetadataDefaults = defaults
	novametadatalog.Info("NovaMetadata defaults initialized", "defaults", defaults)
}

// SetupWebhookWithManager sets up the webhook with the Manager
func (r *NovaMetadata) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-nova-openstack-org-v1beta1-novametadata,mutating=true,failurePolicy=fail,sideEffects=None,groups=nova.openstack.org,resources=novametadata,verbs=create;update,versions=v1beta1,name=mnovametadata.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &NovaMetadata{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *NovaMetadata) Default() {
	novametadatalog.Info("default", "name", r.Name)

	r.Spec.Default()
}

// Default - set defaults for this NovaMetadata spec
func (spec *NovaMetadataSpec) Default() {
	if spec.ContainerImage == "" {
		spec.ContainerImage = novaMetadataDefaults.ContainerImageURL
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-nova-openstack-org-v1beta1-novametadata,mutating=false,failurePolicy=fail,sideEffects=None,groups=nova.openstack.org,resources=novametadata,verbs=create;update,versions=v1beta1,name=vnovametadata.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &NovaMetadata{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *NovaMetadata) ValidateCreate() (admission.Warnings, error) {
	novametadatalog.Info("validate create", "name", r.Name)

	errors := ValidateMetadataDefaultConfigOverwrite(
		field.NewPath("spec").Child("defaultConfigOverwrite"),
		r.Spec.DefaultConfigOverwrite)

	if len(errors) != 0 {
		novametadatalog.Info("validation failed", "name", r.Name)
		return nil, apierrors.NewInvalid(
			schema.GroupKind{Group: "nova.openstack.org", Kind: "NovaMetadata"},
			r.Name, errors)
	}
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *NovaMetadata) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	novametadatalog.Info("validate update", "name", r.Name)

	oldMetadata, ok := old.(*NovaMetadata)
	if !ok || oldMetadata == nil {
		return nil, apierrors.NewInternalError(fmt.Errorf("unable to convert existing object"))
	}

	novametadatalog.Info("validate update", "diff", cmp.Diff(oldMetadata, r))

	errors := ValidateMetadataDefaultConfigOverwrite(
		field.NewPath("spec").Child("defaultConfigOverwrite"),
		r.Spec.DefaultConfigOverwrite)

	if len(errors) != 0 {
		novametadatalog.Info("validation failed", "name", r.Name)
		return nil, apierrors.NewInvalid(
			schema.GroupKind{Group: "nova.openstack.org", Kind: "NovaMetadata"},
			r.Name, errors)
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *NovaMetadata) ValidateDelete() (admission.Warnings, error) {
	novametadatalog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

// ValidateCell0 validates cell0 Metadata template. This is expected to be called
// by higher level validation webhooks
func (r *NovaMetadataTemplate) ValidateCell0(basePath *field.Path) field.ErrorList {
	var errors field.ErrorList
	if *r.Enabled {
		errors = append(
			errors,
			field.Invalid(
				basePath.Child("enabled"), *r.Enabled, "should be false for cell0"),
		)
	}
	return errors
}

func (r *NovaMetadataTemplate) ValidateDefaultConfigOverwrite(basePath *field.Path) field.ErrorList {
	return ValidateMetadataDefaultConfigOverwrite(
		basePath.Child("defaultConfigOverwrite"),
		r.DefaultConfigOverwrite)
}

func ValidateMetadataDefaultConfigOverwrite(
	basePath *field.Path,
	defaultConfigOverwrite map[string]string,
) field.ErrorList {
	return ValidateDefaultConfigOverwrite(
		basePath, defaultConfigOverwrite, []string{"api-paste.ini"})
}
