package groups_v1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"text/template"

	"os"

	nomad "github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/jobspec"
	"github.com/joyent/triton-service-groups/server/handlers"
	"github.com/joyent/triton-service-groups/templates"
)

type OrchestratorJob struct {
	AccountID           int64
	JobName             string
	HealthCheckInterval int
	DesiredCount        int
	InstanceNamePrefix  string
	PackageID           string
	ImageID             string
	ServiceGroupName    string
	UserData            string
	FirewallEnabled     bool
	Networks            []string
	Tags                map[string]string
	MetaData            map[string]string
	TritonAccount       string
	TritonURL           string
	TritonKeyID         string
	TritonKeyMaterial   string
}

func SubmitOrchestratorJob(ctx context.Context, group *ServiceGroup) error {
	session := handlers.GetAuthSession(ctx)

	t, found := templates_v1.FindTemplateByID(ctx, group.TemplateID, session.AccountID)
	if !found {
		return errors.New("Error finding template by ID")
	}

	job, err := prepareJob(t, group)
	if err != nil {
		return err
	}

	deployed, err := registerJob(job)
	if err != nil {
		return err
	}

	log.Print(deployed)

	return nil
}

func UpdateOrchestratorJob(ctx context.Context, group *ServiceGroup) error {
	session := handlers.GetAuthSession(ctx)

	t, found := templates_v1.FindTemplateByID(ctx, group.TemplateID, session.AccountID)
	if !found {
		return errors.New("Error finding template by ID")
	}

	job, err := prepareJob(t, group)
	if err != nil {
		return err
	}

	//we always delete the olb job
	_, err = deregisterJob(*job.ID)
	if err != nil {
		return err
	}

	_, err = registerJob(job)
	if err != nil {
		return err
	}

	return nil
}

func DeleteOrchestratorJob(ctx context.Context, group *ServiceGroup) error {
	session := handlers.GetAuthSession(ctx)

	t, found := templates_v1.FindTemplateByID(ctx, group.TemplateID, session.AccountID)
	if !found {
		return errors.New("Error finding template by ID")
	}

	g := group
	g.Capacity = 0
	job, err := prepareJob(t, g)
	if err != nil {
		return err
	}

	// Delete current version of the job
	_, err = deregisterJob(*job.ID)
	if err != nil {
		return err
	}

	// Submit a new version of the job with a count of 0
	_, err = registerJob(job)
	if err != nil {
		return err
	}

	// Delete current version of the job
	_, err = deregisterJob(*job.ID)
	if err != nil {
		return err
	}

	return nil
}

func deregisterJob(jobID string) (bool, error) {
	orchestratorUrl := os.Getenv("NOMAD_URL")

	client, err := newNomadClient(orchestratorUrl)
	if err != nil {
		return false, err
	}

	_, _, err = client.Jobs().Deregister(jobID, true, nil)
	if err != nil {
		return false, fmt.Errorf("Unable to deregister job with Nomad: %v", err)
	}

	return true, nil
}

func registerJob(job *nomad.Job) (bool, error) {
	orchestratorUrl := os.Getenv("NOMAD_URL")

	client, err := newNomadClient(orchestratorUrl)
	if err != nil {
		return false, err
	}

	_, _, err = client.Jobs().Validate(job, nil)
	if err != nil {
		return false, fmt.Errorf("Failed to validate Nomad Job: %v", err)
	}

	_, _, err = client.Jobs().Register(job, nil)
	if err != nil {
		return false, fmt.Errorf("Unable to register job with Nomad: %v", err)
	}

	_, _, err = client.Jobs().PeriodicForce(*job.ID, nil)
	if err != nil {
		return false, fmt.Errorf("Unable to trigger a periodic instance of job: %v", err)
	}

	return true, nil
}

func prepareJob(t *templates_v1.InstanceTemplate, group *ServiceGroup) (*nomad.Job, error) {
	tpl := &bytes.Buffer{}
	details := createJobDetails(t, group)

	fmap := template.FuncMap{
		"formatAsMinutes": formatAsMinutes,
	}

	jobT := template.Must(template.New("job").Funcs(fmap).Parse(jobTemplate))
	err := jobT.Execute(tpl, details)
	if err != nil {
		return nil, err
	}

	job, err := jobspec.Parse(tpl)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func createJobDetails(template *templates_v1.InstanceTemplate, group *ServiceGroup) OrchestratorJob {
	job := OrchestratorJob{
		AccountID:           group.AccountID,
		JobName:             fmt.Sprintf("%s_%d", group.GroupName, template.ID),
		HealthCheckInterval: group.HealthCheckInterval,
		DesiredCount:        group.Capacity,
		PackageID:           template.Package,
		ImageID:             template.ImageID,
		ServiceGroupName:    group.GroupName,
		FirewallEnabled:     template.FirewallEnabled,
	}

	if template.InstanceNamePrefix != "" {
		job.InstanceNamePrefix = template.InstanceNamePrefix
	}

	if template.UserData != "" {
		job.UserData = template.UserData
	}

	if len(template.Networks) > 0 {
		job.Networks = template.Networks
	}

	if template.Tags != nil {
		job.Tags = template.Tags
	}

	if template.MetaData != nil {
		job.MetaData = template.MetaData
	}

	job.TritonAccount = os.Getenv("TRITON_ACCOUNT")
	job.TritonURL = os.Getenv("TRITON_URL")
	job.TritonKeyID = os.Getenv("TRITON_KEY_ID")

	keyMaterial := strings.Replace(os.Getenv("TRITON_KEY_MATERIAL"), "\n", `\n`, -1)
	job.TritonKeyMaterial = keyMaterial

	return job
}

func formatAsMinutes(interval int) (string, error) {
	return fmt.Sprintf("%d", interval/60), nil
}

const jobTemplate = `
job "{{.JobName}}" {
  type = "batch"
  periodic {
	cron = "*/{{.HealthCheckInterval | formatAsMinutes}} * * * * "
	prohibit_overlap = true
  }
  datacenters = ["dc1"]
  group "scale" {
	task "healthy" {
	  driver = "exec"
	  artifact {
		source = "http://us-east.manta.joyent.com/productci/public/tsg"
	  }
	  config {
		command = "tsg"
		args = [
		  "scale",
		  "--count", "{{ .DesiredCount }}",
		  "--pkg-id", "{{ .PackageID }}",
		  "--img-id", "{{ .ImageID }}",
		  "--tsg-name", "{{ .ServiceGroupName }}",
		  {{if .InstanceNamePrefix -}}
		  "--name-prefix", "{{ .InstanceNamePrefix }}",
		  {{- end }}
		  {{if .UserData -}}
		  "--userdata", "{{ .UserData }}",
		  {{- end }}
		  {{range .Networks}}
		  "--networks", "{{ . }}",
		  {{- end }}
		  {{range $key, $value := .Tags}}
		  "--tag", "{{$key}}={{$value}}",
		  {{- end }}
		  {{range $key, $value := .MetaData}}
		  "--metadata", "{{$key}}={{$value}}",
		  {{- end }}
		  "-A", "{{ .TritonAccount }}",
		  "-K", "{{ .TritonKeyID }}",
		  "-U", "{{ .TritonURL }}",
		  {{if .TritonKeyMaterial -}}
		  "--key-material", "{{ .TritonKeyMaterial }}",
		  {{- end}}
		]
	  }
	}
  }
}
`

func newNomadClient(addr string) (*nomad.Client, error) {
	config := nomad.DefaultConfig()

	if addr != "" {
		config.Address = addr
	}

	c, err := nomad.NewClient(config)
	if err != nil {
		return nil, err
	}

	return c, nil
}
