package validate

import (
	"github.com/xcloudnative/xcloud/pkg/apis/xcloudnative.io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/argoproj/pkg/json"
	wfv1 "github.com/xcloudnative/xcloud/pkg/apis/xcloudnative.io/v1alpha1"
	"github.com/xcloudnative/xcloud/pkg/errors"
	"github.com/xcloudnative/xcloud/pkg/xflow/common"
)

// LintWorkflowDir validates all workflow manifests in a directory. Ignores non-workflow manifests
func LintWorkflowDir(dirPath string, strict bool) error {
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		fileExt := filepath.Ext(info.Name())
		switch fileExt {
		case ".yaml", ".yml", ".json":
		default:
			return nil
		}
		return LintWorkflowFile(path, strict)
	}
	return filepath.Walk(dirPath, walkFunc)
}

// LintWorkflowFile lints a json file, or multiple workflow manifest in a single yaml file. Ignores
// non-workflow manifests
func LintWorkflowFile(filePath string, strict bool) error {
	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		return errors.Errorf(errors.CodeBadRequest, "Can't read from file: %s, err: %v", filePath, err)
	}
	var workflows []wfv1.Xflow
	if json.IsJSON(body) {
		var wf wfv1.Xflow
		if strict {
			err = json.UnmarshalStrict(body, &wf)
		} else {
			err = json.Unmarshal(body, &wf)
		}
		if err == nil {
			workflows = []wfv1.Xflow{wf}
		} else {
			if wf.Kind != "" && wf.Kind != xcloudnativeio.Kind {
				// If we get here, it was a k8s manifest which was not of type 'Workflow'
				// We ignore these since we only care about validating Workflow manifests.
				return nil
			}
		}
	} else {
		workflows, err = common.SplitYAMLFile(body, strict)
	}
	if err != nil {
		return errors.Errorf(errors.CodeBadRequest, "%s failed to parse: %v", filePath, err)
	}
	for _, wf := range workflows {
		err = ValidateWorkflow(&wf, true)
		if err != nil {
			return errors.Errorf(errors.CodeBadRequest, "%s: %s", filePath, err.Error())
		}
	}
	return nil
}
