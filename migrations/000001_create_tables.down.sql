-- Drop dependent tables first (leaf nodes)
DROP TABLE IF EXISTS schedules;
DROP TABLE IF EXISTS user_groups;
DROP TABLE IF EXISTS user_subscriptions;

-- Drop tables referencing other tables
DROP TABLE IF EXISTS training_groups;
DROP TABLE IF EXISTS trainers;

-- Drop core tables with dependencies
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;

-- Drop category and location tables
DROP TABLE IF EXISTS group_category;
DROP TABLE IF EXISTS pools;

-- Drop subscription-related tables (after user_subscriptions)
DROP TABLE IF EXISTS subscriptions;