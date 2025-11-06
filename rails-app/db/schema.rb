# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# This file is the source Rails uses to define your schema when running `bin/rails
# db:schema:load`. When creating a new database, `bin/rails db:schema:load` tends to
# be faster and is potentially less error prone than running all of your
# migrations from scratch. Old migrations may fail to apply correctly if those
# migrations use external dependencies or application code.
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema[7.1].define(version: 2024_01_01_000007) do
  # These are extensions that must be enabled in order to support this database
  enable_extension "plpgsql"

  create_table "telegram_accounts", force: :cascade do |t|
    t.string "phone_number", null: false
    t.string "session_data", limit: 10000
    t.string "status", default: "pending", null: false
    t.string "first_name"
    t.string "last_name"
    t.string "username"
    t.bigint "telegram_user_id"
    t.datetime "last_active_at"
    t.text "error_message"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["phone_number"], name: "index_telegram_accounts_on_phone_number", unique: true
    t.index ["status"], name: "index_telegram_accounts_on_status"
    t.index ["telegram_user_id"], name: "index_telegram_accounts_on_telegram_user_id", unique: true
  end

  create_table "message_templates", force: :cascade do |t|
    t.string "name", null: false
    t.text "content", null: false
    t.string "media_type"
    t.string "media_url"
    t.boolean "parse_mode_enabled", default: false
    t.string "parse_mode"
    t.boolean "disable_web_page_preview", default: false
    t.json "buttons"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["name"], name: "index_message_templates_on_name"
  end

  create_table "channels", force: :cascade do |t|
    t.string "name", null: false
    t.string "channel_type", null: false
    t.bigint "telegram_id", null: false
    t.string "username"
    t.string "title"
    t.integer "members_count"
    t.text "description"
    t.datetime "last_synced_at"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["telegram_id"], name: "index_channels_on_telegram_id", unique: true
    t.index ["username"], name: "index_channels_on_username"
    t.index ["channel_type"], name: "index_channels_on_channel_type"
  end

  create_table "schedules", force: :cascade do |t|
    t.string "name", null: false
    t.bigint "telegram_account_id", null: false
    t.bigint "message_template_id", null: false
    t.string "schedule_type", null: false
    t.string "cron_expression"
    t.integer "interval_minutes"
    t.datetime "scheduled_at"
    t.string "timezone", default: "UTC"
    t.boolean "active", default: false
    t.datetime "next_run_at"
    t.datetime "last_run_at"
    t.integer "total_runs", default: 0
    t.integer "successful_runs", default: 0
    t.integer "failed_runs", default: 0
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["telegram_account_id"], name: "index_schedules_on_telegram_account_id"
    t.index ["message_template_id"], name: "index_schedules_on_message_template_id"
    t.index ["active"], name: "index_schedules_on_active"
    t.index ["next_run_at"], name: "index_schedules_on_next_run_at"
    t.index ["schedule_type"], name: "index_schedules_on_schedule_type"
  end

  create_table "schedule_channels", force: :cascade do |t|
    t.bigint "schedule_id", null: false
    t.bigint "channel_id", null: false
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["schedule_id"], name: "index_schedule_channels_on_schedule_id"
    t.index ["channel_id"], name: "index_schedule_channels_on_channel_id"
    t.index ["schedule_id", "channel_id"], name: "index_schedule_channels_on_schedule_id_and_channel_id", unique: true
  end

  create_table "send_logs", force: :cascade do |t|
    t.bigint "schedule_id", null: false
    t.bigint "telegram_account_id", null: false
    t.bigint "channel_id", null: false
    t.string "status", null: false
    t.text "message_content"
    t.bigint "telegram_message_id"
    t.text "error_message"
    t.datetime "sent_at"
    t.integer "retry_count", default: 0
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["schedule_id"], name: "index_send_logs_on_schedule_id"
    t.index ["telegram_account_id"], name: "index_send_logs_on_telegram_account_id"
    t.index ["channel_id"], name: "index_send_logs_on_channel_id"
    t.index ["status"], name: "index_send_logs_on_status"
    t.index ["sent_at"], name: "index_send_logs_on_sent_at"
    t.index ["created_at"], name: "index_send_logs_on_created_at"
  end

  create_table "users", force: :cascade do |t|
    t.string "email", null: false
    t.string "password_digest", null: false
    t.string "role", default: "admin", null: false
    t.datetime "last_login_at"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["email"], name: "index_users_on_email", unique: true
  end

  add_foreign_key "schedules", "telegram_accounts"
  add_foreign_key "schedules", "message_templates"
  add_foreign_key "schedule_channels", "schedules"
  add_foreign_key "schedule_channels", "channels"
  add_foreign_key "send_logs", "schedules"
  add_foreign_key "send_logs", "telegram_accounts"
  add_foreign_key "send_logs", "channels"
end
