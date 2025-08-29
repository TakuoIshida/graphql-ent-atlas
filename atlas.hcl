// Atlas configuration file for Ent-based database migrations

// Define the environments
env "local" {
  src = "ent://ent/schema"
  dev = "docker://postgres/15/dev?search_path=public"
  url = getenv("DATABASE_URL") != "" ? getenv("DATABASE_URL") : "postgres://postgres:password@localhost:5432/todoapp?sslmode=disable"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "prod" {
  src = "ent://ent/schema"
  url = getenv("DATABASE_URL")
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}