job "tsg-v1" {
  type = "service"

  datacenters = [
    "dc1"
  ]

  constraint {
    attribute = "${attr.kernel.name}"
    value     = "linux"
  }

  update {
    health_check = "task_states"
    max_parallel = 1
    stagger      = "10s"
  }

  group "deployment" {
    count = 1

    constraint {
      operator = "distinct_hosts"
      value    = "true"
    }

    task "api" {
      driver = "exec"

      artifact {
        source = "https://github.com/joyent/triton-service-groups/releases/download/v0.1.0/triton-service-groups_0.1.0_linux_amd64.tar.gz"
      }

      # env {
      #   "TSG_CRDB_HOST" = "172.27.10.11"
      # }

      config {
        command = "triton-sg"
        args = [
          "agent",
          "--log-level", "DEBUG"
        ]
      }

      service {
        port = "http"
        tags = [
          "api"
        ]
      }

      resources {
        network {
          port "http" {}
        }
      }
    }
  }
}
