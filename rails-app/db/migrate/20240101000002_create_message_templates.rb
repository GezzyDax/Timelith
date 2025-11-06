class CreateMessageTemplates < ActiveRecord::Migration[7.1]
  def change
    create_table :message_templates do |t|
      t.string :name, null: false
      t.text :content, null: false
      t.string :media_type # photo, video, document, etc.
      t.string :media_url
      t.boolean :parse_mode_enabled, default: false
      t.string :parse_mode # markdown, html
      t.boolean :disable_web_page_preview, default: false
      t.json :buttons # for inline keyboards
      t.timestamps
    end

    add_index :message_templates, :name
  end
end
