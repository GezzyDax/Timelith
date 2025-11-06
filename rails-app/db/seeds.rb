# This file should contain all the record creation needed to seed the database with its default values.
# The data can then be loaded with the bin/rails db:seed command (or created alongside the database with db:setup).

# Create default admin user
if User.count.zero?
  User.create!(
    email: 'admin@example.com',
    password: 'admin123',
    role: 'admin'
  )
  puts "Created admin user: admin@example.com / admin123"
  puts "⚠️  IMPORTANT: Change this password after first login!"
end
