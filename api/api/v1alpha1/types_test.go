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

package v1alpha1_test

import (
	"os"
	"testing"

	"github.com/kluster-management/fluxcd-addon/api/api/v1alpha1"

	sc "kmodules.xyz/schema-checker"
)

func TestDefaultValues(t *testing.T) {
	checker := sc.New(os.DirFS("../../.."),
		sc.TestCase{Obj: v1alpha1.FluxCDConfigSpec{}, File: "https://github.com/fluxcd-community/helm-charts/raw/flux2-2.10.3/charts/flux2/values.yaml"},
	)
	checker.TestAll(t)
}
