package groups_v1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	stdlog "log"
	"strings"
	"text/template"

	nomad "github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/jobspec"
	"github.com/joyent/triton-service-groups/accounts"
	"github.com/joyent/triton-service-groups/server/handlers"
	"github.com/joyent/triton-service-groups/templates"
	"github.com/rs/zerolog/log"
)

type OrchestratorJob struct {
	Datacenter        string
	JobName           string
	DesiredCount      int
	PackageID         string
	ImageID           string
	ServiceGroupName  string
	TemplateID        string
	UserData          string
	FirewallEnabled   bool
	Networks          []string
	Tags              map[string]string
	MetaData          map[string]string
	TritonAccount     string
	TritonURL         string
	TritonKeyID       string
	TritonKeyMaterial string
}

func SubmitOrchestratorJob(ctx context.Context, group *ServiceGroup) error {
	session := handlers.GetAuthSession(ctx)

	t, found := templates_v1.FindTemplateByID(ctx, group.TemplateID, session.AccountID)
	if !found {
		return errors.New("Error finding template by ID")
	}

	job, err := prepareJob(ctx, t, group)
	if err != nil {
		return err
	}

	deployed, err := registerJob(ctx, job)
	if err != nil {
		return err
	}

	stdlog.Print(deployed)

	return nil
}

func UpdateOrchestratorJob(ctx context.Context, group *ServiceGroup) error {
	session := handlers.GetAuthSession(ctx)

	t, found := templates_v1.FindTemplateByID(ctx, group.TemplateID, session.AccountID)
	if !found {
		return errors.New("Error finding template by ID")
	}

	job, err := prepareJob(ctx, t, group)
	if err != nil {
		return err
	}

	// we always delete the old job
	_, err = deregisterJob(ctx, *job.ID)
	if err != nil {
		return err
	}

	_, err = registerJob(ctx, job)
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
	job, err := prepareJob(ctx, t, g)
	if err != nil {
		return err
	}

	// Delete current version of the job
	_, err = deregisterJob(ctx, *job.ID)
	if err != nil {
		return err
	}

	// Submit a new version of the job with a count of 0
	_, err = registerJob(ctx, job)
	if err != nil {
		return err
	}

	// Delete current version of the job
	_, err = deregisterJob(ctx, *job.ID)
	if err != nil {
		return err
	}

	return nil
}

func deregisterJob(ctx context.Context, jobID string) (bool, error) {
	client, ok := handlers.GetNomadClient(ctx)
	if !ok {
		log.Fatal().Err(handlers.ErrNoNomadClient)
		return false, handlers.ErrNoNomadClient
	}

	_, _, err := client.Jobs().Deregister(jobID, true, nil)
	if err != nil {
		return false, fmt.Errorf("Unable to deregister job with Nomad: %v", err)
	}

	return true, nil
}

func registerJob(ctx context.Context, job *nomad.Job) (bool, error) {
	client, ok := handlers.GetNomadClient(ctx)
	if !ok {
		log.Error().Err(handlers.ErrNoNomadClient)
		return false, handlers.ErrNoNomadClient
	}

	_, _, err := client.Jobs().Validate(job, nil)
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

func prepareJob(ctx context.Context, t *templates_v1.InstanceTemplate, group *ServiceGroup) (*nomad.Job, error) {
	session := handlers.GetAuthSession(ctx)

	tpl := &bytes.Buffer{}
	details := createJobDetails(t, group)
	details.Datacenter = session.Datacenter
	if err := details.getTritonAccountDetails(ctx); err != nil {
		return nil, err
	}

	jobT := template.Must(template.New("job").Parse(jobTemplate))
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

func (j *OrchestratorJob) getTritonAccountDetails(ctx context.Context) error {
	session := handlers.GetAuthSession(ctx)

	db, ok := handlers.GetDBPool(ctx)
	if !ok {
		log.Error().Err(handlers.ErrNoConnPool)
		return handlers.ErrNoConnPool
	}

	store := accounts.NewStore(db)

	account, err := store.FindByID(ctx, session.AccountID)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	credential, err := account.GetTritonCredential(ctx)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	log.Debug().
		Str("account_id", account.ID).
		Str("account_name", account.AccountName).
		Str("fingerprint", credential.KeyID).
		Msg("orchestrator: found triton credentials for account")

	j.TritonKeyMaterial = strings.Replace(credential.KeyMaterial, "\n", `\n`, -1)
	j.TritonAccount = credential.AccountName
	j.TritonKeyID = credential.KeyID
	j.TritonURL = session.TritonURL

	j.JobName = fmt.Sprintf("%s_%s", j.ServiceGroupName, account.TritonUUID)

	return nil
}

func createJobDetails(template *templates_v1.InstanceTemplate, group *ServiceGroup) OrchestratorJob {
	job := OrchestratorJob{
		DesiredCount:     group.Capacity,
		PackageID:        template.Package,
		ImageID:          template.ImageID,
		ServiceGroupName: group.GroupName,
		FirewallEnabled:  template.FirewallEnabled,
		TemplateID:       template.ID,
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

const jobTemplate = `
job "{{.JobName}}" {
  type = "batch"
  periodic {
	cron = "*/2 * * * * *"
	prohibit_overlap = true
  }
  datacenters = ["{{ .Datacenter }}"]
  group "scale" {
	constraint {
	  distinct_hosts = true
	}
	constraint {
	  operator = "="
	  attribute = "${meta.role}"
	  value = "automater"
	}
	task "healthy" {
	  driver = "exec"
	  artifact {
		source = "https://github.com/joyent/tsg-cli/releases/download/v0.1.0/tsg-cli_0.1.0_linux_amd64.tar.gz"
		options {
		  checksum = "sha256:0aa5bf0bfea4da8b4c58178e0764c0fa06ca09075df3a44ef6ed6fc11d2b74ce"
		}
	  }
	  config {
		command = "tsg-cli"
		args = [
		  "scale",
		  "--count", "{{ .DesiredCount }}",
		  "--pkg-id", "{{ .PackageID }}",
		  "--img-id", "{{ .ImageID }}",
		  "--tsg-name", "{{ .ServiceGroupName }}",
		  "--template-id", "{{ .TemplateID }}",
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
