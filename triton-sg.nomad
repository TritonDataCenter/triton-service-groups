job "tsg-v1" {

  type = "service"

  datacenters = ["dc1"]

  group "deployment" {

    constraint {
      distinct_hosts = true
    }

    # constraint {
    #   operator = "distinct_property"
    #   attribute = "${meta.role}"
    #   value = "api-client"
    # }

    count = 1

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
        tags = ["api"]
        port = "http"
      }

      resources {
        network {
          port "http" {}
        }
      }
    }
  }
}
