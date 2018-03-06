package groups_v1

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"text/template"

	nomad "github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/jobspec"
	"github.com/joyent/triton-service-groups/session"
	"github.com/joyent/triton-service-groups/templates"
	"github.com/y0ssar1an/q"
)

type OrchestratorJob struct {
	AccountID           string
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
}

func SubmitOrchestratorJob(session *session.TsgSession, group *ServiceGroup) error {

	t, found := templates_v1.FindTemplateByID(session.DbPool, group.TemplateId, session.AccountId)
	if !found {
		return errors.New("Error finding template by ID")
	}

	job, err := prepareJob(t, group)
	if err != nil {
		return err
	}

	q.Q(job)

	deployed, err := deployJob(job)
	if err != nil {
		return err
	}

	log.Print(deployed)

	return nil
}

func UpdateOrchestratorJob(session *session.TsgSession, group *ServiceGroup) error {
	t, found := templates_v1.FindTemplateByID(session.DbPool, group.TemplateId, session.AccountId)
	if !found {
		return errors.New("Error finding template by ID")
	}

	job, err := prepareJob(t, group)
	if err != nil {
		return err
	}

	//we always delete the olb job
	_, err = deleteJob(*job.ID)
	if err != nil {
		return err
	}

	_, err = deployJob(job)
	if err != nil {
		return err
	}

	return nil
}

func DeleteOrchestratorJob(session *session.TsgSession, group *ServiceGroup) error {
	t, found := templates_v1.FindTemplateByID(session.DbPool, group.TemplateId, session.AccountId)
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
	_, err = deleteJob(*job.ID)
	if err != nil {
		return err
	}

	// Submit a new version of the job with a count of 0
	_, err = deployJob(job)
	if err != nil {
		return err
	}

	// Delete current version of the job
	_, err = deleteJob(*job.ID)
	if err != nil {
		return err
	}

	return nil
}

func deleteJob(jobID string) (bool, error) {

	client, err := newNomadClient("http://localhost:4646")
	if err != nil {
		return false, err
	}

	_, _, err = client.Jobs().Deregister(jobID, true, nil)
	if err != nil {
		return false, fmt.Errorf("Unable to deregister job with Nomad: %v", err)
	}

	return true, nil
}

func deployJob(job *nomad.Job) (bool, error) {

	client, err := newNomadClient("http://localhost:4646")
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
		panic(err)
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
		AccountID:           group.AccountId,
		JobName:             fmt.Sprintf("%s_%d", group.GroupName, template.ID),
		HealthCheckInterval: group.HealthCheckInterval,
		DesiredCount:        group.Capacity,
		PackageID:           template.Package,
		ImageID:             template.ImageId,
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
          "-A", "stack72_joyent",
          "-K", "aa:a1:50:b0:07:52:7b:6b:3c:ad:11:e1:dd:41:7e:08",
          "-U", "https://us-sw-1.api.joyentcloud.com",
          "--key-material", "-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
"
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
