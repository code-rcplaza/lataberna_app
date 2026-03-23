// Atlas HCL schema — Module 0: Magic Link Authentication
// SQLite stores datetimes as TEXT in ISO 8601 (RFC3339) format.

table "users" {
  schema = schema.main

  column "id" {
    type    = text
    null    = false
  }
  column "email" {
    type    = text
    null    = false
  }
  column "created_at" {
    type    = text
    null    = false
  }

  primary_key {
    columns = [column.id]
  }

  index "users_email_unique" {
    columns = [column.email]
    unique  = true
  }
}

table "magic_link_tokens" {
  schema = schema.main

  column "id" {
    type    = text
    null    = false
  }
  column "hashed_token" {
    type    = text
    null    = false
  }
  column "email" {
    type    = text
    null    = false
  }
  column "expires_at" {
    type    = text
    null    = false
  }
  column "used_at" {
    type    = text
    null    = true
  }
  column "created_at" {
    type    = text
    null    = false
  }

  primary_key {
    columns = [column.id]
  }

  index "magic_link_tokens_hashed_token_unique" {
    columns = [column.hashed_token]
    unique  = true
  }
}

table "sessions" {
  schema = schema.main

  column "id" {
    type    = text
    null    = false
  }
  column "user_id" {
    type    = text
    null    = false
  }
  column "expires_at" {
    type    = text
    null    = false
  }
  column "created_at" {
    type    = text
    null    = false
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "sessions_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }
}

table "characters" {
  schema = schema.main

  column "id" {
    type = text
    null = false
  }
  column "user_id" {
    type = text
    null = false
  }
  column "name" {
    type = text
    null = false
  }
  column "species" {
    type = text
    null = false
  }
  column "sub_species" {
    type = text
    null = true
  }
  column "class" {
    type = text
    null = false
  }
  column "level" {
    type    = integer
    null    = false
  }
  column "ruleset" {
    type = text
    null = false
  }
  column "ability_bonus_source" {
    type = text
    null = false
  }
  column "base_stats" {
    type = text
    null = false
  }
  column "final_stats" {
    type = text
    null = false
  }
  column "modifiers" {
    type = text
    null = false
  }
  column "derived" {
    type = text
    null = false
  }
  column "background" {
    type = text
    null = false
  }
  column "motivation" {
    type = text
    null = false
  }
  column "secret" {
    type = text
    null = false
  }
  column "locks" {
    type = text
    null = false
  }
  column "seed" {
    type = integer
    null = true
  }
  column "created_at" {
    type = text
    null = false
  }
  column "updated_at" {
    type = text
    null = false
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "characters_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete   = NO_ACTION
    on_update   = NO_ACTION
  }

  index "characters_user_id_idx" {
    columns = [column.user_id]
  }
}

schema "main" {}
