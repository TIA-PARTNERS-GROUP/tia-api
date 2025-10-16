-- CreateTable
CREATE TABLE `business_connections` (
    `id` INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
    `initiating_business_id` INTEGER UNSIGNED NOT NULL,
    `receiving_business_id` INTEGER UNSIGNED NOT NULL,
    `connection_type` ENUM('Partnership', 'Supplier', 'Client', 'Referral', 'Collaboration') NOT NULL,
    `status` ENUM('pending', 'active', 'rejected', 'inactive') NOT NULL DEFAULT 'pending',
    `initiated_by_user_id` INTEGER UNSIGNED NOT NULL,
    `notes` TEXT NULL,
    `created_at` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),
    `updated_at` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),

    INDEX `fk_business_connections_initiated_by`(`initiated_by_user_id`),
    INDEX `idx_business_connections_initiating_id`(`initiating_business_id`),
    INDEX `idx_business_connections_receiving_id`(`receiving_business_id`),
    INDEX `idx_business_connections_status`(`status`),
    UNIQUE INDEX `uq_business_connections_unique`(`initiating_business_id`, `receiving_business_id`, `connection_type`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `business_tags` (
    `id` INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
    `business_id` INTEGER UNSIGNED NOT NULL,
    `tag_type` ENUM('client', 'service', 'specialty') NOT NULL,
    `description` VARCHAR(100) NOT NULL,
    `created_at` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),

    INDEX `idx_business_tags_business_id`(`business_id`),
    INDEX `idx_business_tags_tag_type`(`tag_type`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `businesses` (
    `id` INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
    `operator_user_id` INTEGER UNSIGNED NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `tagline` VARCHAR(100) NULL,
    `website` VARCHAR(255) NULL,
    `contact_name` VARCHAR(60) NULL,
    `contact_phone_no` VARCHAR(20) NULL,
    `contact_email` VARCHAR(254) NULL,
    `description` TEXT NULL,
    `address` VARCHAR(100) NULL,
    `city` VARCHAR(60) NULL,
    `state` VARCHAR(60) NULL,
    `country` VARCHAR(60) NULL,
    `postal_code` VARCHAR(20) NULL,
    `value` DECIMAL(15, 2) NULL,
    `business_type` ENUM('Consulting', 'Retail', 'Technology', 'Manufacturing', 'Services', 'Other') NOT NULL,
    `business_category` ENUM('B2B', 'B2C', 'Non-Profit', 'Government', 'Mixed') NOT NULL,
    `business_phase` ENUM('Startup', 'Growth', 'Mature', 'Exit') NOT NULL,
    `active` TINYINT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),
    `updated_at` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),

    INDEX `idx_businesses_active`(`active`),
    INDEX `idx_businesses_business_type`(`business_type`),
    INDEX `idx_businesses_operator_user_id`(`operator_user_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `idea_votes` (
    `voter_user_id` INTEGER UNSIGNED NOT NULL,
    `idea_id` INTEGER UNSIGNED NOT NULL,
    `vote_type` ENUM('up', 'down') NOT NULL,
    `created_at` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),

    INDEX `idx_idea_votes_idea_id`(`idea_id`),
    PRIMARY KEY (`voter_user_id`, `idea_id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `ideas` (
    `id` INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
    `submitted_by_user_id` INTEGER UNSIGNED NOT NULL,
    `title` VARCHAR(200) NOT NULL,
    `content` TEXT NOT NULL,
    `status` ENUM('open', 'under_review', 'planned', 'in_progress', 'completed', 'rejected') NOT NULL DEFAULT 'open',
    `created_at` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),
    `updated_at` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),

    INDEX `idx_ideas_status`(`status`),
    INDEX `idx_ideas_user_id`(`submitted_by_user_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `notifications` (
    `id` INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
    `sender_user_id` INTEGER UNSIGNED NULL,
    `receiver_user_id` INTEGER UNSIGNED NOT NULL,
    `notification_type` ENUM('connection_request', 'project_invite', 'message', 'system') NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `message` TEXT NOT NULL,
    `related_entity_type` ENUM('business', 'project', 'publication', 'idea') NULL,
    `related_entity_id` INTEGER UNSIGNED NULL,
    `read` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `action_url` VARCHAR(500) NULL,
    `created_at` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),

    INDEX `idx_notifications_created_at`(`created_at`),
    INDEX `idx_notifications_read`(`read`),
    INDEX `idx_notifications_receiver_user_id`(`receiver_user_id`),
    INDEX `idx_notifications_sender_user_id`(`sender_user_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `project_members` (
    `project_id` INTEGER UNSIGNED NOT NULL,
    `user_id` INTEGER UNSIGNED NOT NULL,
    `role` ENUM('manager', 'contributor', 'reviewer') NOT NULL DEFAULT 'contributor',
    `joined_at` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),

    INDEX `idx_project_members_user_id`(`user_id`),
    PRIMARY KEY (`project_id`, `user_id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `project_skills` (
    `project_id` INTEGER UNSIGNED NOT NULL,
    `skill_id` INTEGER UNSIGNED NOT NULL,
    `importance` ENUM('required', 'preferred', 'optional') NOT NULL DEFAULT 'preferred',

    INDEX `idx_project_skills_skill_id`(`skill_id`),
    PRIMARY KEY (`project_id`, `skill_id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `projects` (
    `id` INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
    `managed_by_user_id` INTEGER UNSIGNED NOT NULL,
    `business_id` INTEGER UNSIGNED NULL,
    `name` VARCHAR(100) NOT NULL,
    `description` TEXT NULL,
    `project_status` ENUM('planning', 'active', 'on_hold', 'completed', 'cancelled') NOT NULL DEFAULT 'planning',
    `start_date` DATE NULL,
    `target_end_date` DATE NULL,
    `actual_end_date` DATE NULL,
    `created_at` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),
    `updated_at` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),

    INDEX `idx_projects_business_id`(`business_id`),
    INDEX `idx_projects_managed_by_user_id`(`managed_by_user_id`),
    INDEX `idx_projects_status`(`project_status`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `publications` (
    `id` INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` INTEGER UNSIGNED NOT NULL,
    `business_id` INTEGER UNSIGNED NULL,
    `publication_type` ENUM('post', 'case_study', 'testimonial', 'article') NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `slug` VARCHAR(300) NOT NULL,
    `excerpt` TEXT NULL,
    `content` LONGTEXT NOT NULL,
    `thumbnail` VARCHAR(255) NULL,
    `video_url` VARCHAR(255) NULL,
    `published` TINYINT NOT NULL DEFAULT 0,
    `published_at` TIMESTAMP(0) NULL,
    `created_at` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),
    `updated_at` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),

    UNIQUE INDEX `slug`(`slug`),
    INDEX `idx_publications_business_id`(`business_id`),
    INDEX `idx_publications_published`(`published`),
    INDEX `idx_publications_type`(`publication_type`),
    INDEX `idx_publications_user_id`(`user_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `skills` (
    `id` INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
    `category` VARCHAR(100) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `description` TEXT NULL,
    `active` TINYINT NOT NULL DEFAULT 1,
    `created_at` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),

    UNIQUE INDEX `uq_skills_name`(`name`),
    INDEX `idx_skills_active`(`active`),
    INDEX `idx_skills_category`(`category`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `user_sessions` (
    `id` INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` INTEGER UNSIGNED NOT NULL,
    `token_hash` CHAR(128) NOT NULL,
    `ip_address` VARCHAR(45) NULL,
    `user_agent` TEXT NULL,
    `created_at` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),
    `expires_at` TIMESTAMP(0) NOT NULL,
    `revoked_at` TIMESTAMP(0) NULL,

    UNIQUE INDEX `uq_user_sessions_token_hash`(`token_hash`),
    INDEX `idx_user_sessions_expires_at`(`expires_at`),
    INDEX `idx_user_sessions_user_id`(`user_id`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `user_skills` (
    `skill_id` INTEGER UNSIGNED NOT NULL,
    `user_id` INTEGER UNSIGNED NOT NULL,
    `proficiency_level` ENUM('beginner', 'intermediate', 'advanced', 'expert') NULL DEFAULT 'intermediate',
    `created_at` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),

    INDEX `idx_user_skills_user_id`(`user_id`),
    PRIMARY KEY (`skill_id`, `user_id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- CreateTable
CREATE TABLE `users` (
    `id` INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
    `first_name` VARCHAR(60) NOT NULL,
    `last_name` VARCHAR(60) NULL,
    `login_email` VARCHAR(254) NOT NULL,
    `password_hash` VARCHAR(255) NULL,
    `contact_email` VARCHAR(254) NULL,
    `contact_phone_no` VARCHAR(20) NULL,
    `adk_session_id` VARCHAR(128) NULL,
    `password_reset_token` BLOB NULL,
    `password_reset_requested_at` TIMESTAMP(0) NULL,
    `email_verified` BOOLEAN NOT NULL DEFAULT false,
    `active` BOOLEAN NOT NULL DEFAULT true,
    `created_at` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),
    `updated_at` TIMESTAMP(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),

    UNIQUE INDEX `login_email`(`login_email`),
    INDEX `idx_users_active`(`active`),
    INDEX `idx_users_contact_email`(`contact_email`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- AddForeignKey
ALTER TABLE `business_connections` ADD CONSTRAINT `fk_business_connections_initiated_by` FOREIGN KEY (`initiated_by_user_id`) REFERENCES `users`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `business_connections` ADD CONSTRAINT `fk_business_connections_initiating_id` FOREIGN KEY (`initiating_business_id`) REFERENCES `businesses`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `business_connections` ADD CONSTRAINT `fk_business_connections_receiving_id` FOREIGN KEY (`receiving_business_id`) REFERENCES `businesses`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `business_tags` ADD CONSTRAINT `fk_business_tags_business_id` FOREIGN KEY (`business_id`) REFERENCES `businesses`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `businesses` ADD CONSTRAINT `fk_businesses_operator_user_id` FOREIGN KEY (`operator_user_id`) REFERENCES `users`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `idea_votes` ADD CONSTRAINT `fk_idea_votes_idea_id` FOREIGN KEY (`idea_id`) REFERENCES `ideas`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `idea_votes` ADD CONSTRAINT `fk_idea_votes_user_id` FOREIGN KEY (`voter_user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `ideas` ADD CONSTRAINT `fk_ideas_user_id` FOREIGN KEY (`submitted_by_user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `notifications` ADD CONSTRAINT `fk_notifications_receiver_user_id` FOREIGN KEY (`receiver_user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `notifications` ADD CONSTRAINT `fk_notifications_sender_user_id` FOREIGN KEY (`sender_user_id`) REFERENCES `users`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `project_members` ADD CONSTRAINT `fk_project_members_project_id` FOREIGN KEY (`project_id`) REFERENCES `projects`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `project_members` ADD CONSTRAINT `fk_project_members_user_id` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `project_skills` ADD CONSTRAINT `fk_project_skills_project_id` FOREIGN KEY (`project_id`) REFERENCES `projects`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `project_skills` ADD CONSTRAINT `fk_project_skills_skill_id` FOREIGN KEY (`skill_id`) REFERENCES `skills`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `projects` ADD CONSTRAINT `fk_projects_business_id` FOREIGN KEY (`business_id`) REFERENCES `businesses`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `projects` ADD CONSTRAINT `fk_projects_managed_by_user_id` FOREIGN KEY (`managed_by_user_id`) REFERENCES `users`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `publications` ADD CONSTRAINT `fk_publications_business_id` FOREIGN KEY (`business_id`) REFERENCES `businesses`(`id`) ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `publications` ADD CONSTRAINT `fk_publications_user_id` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `user_sessions` ADD CONSTRAINT `fk_user_sessions_user_id` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `user_skills` ADD CONSTRAINT `fk_user_skills_skill_id` FOREIGN KEY (`skill_id`) REFERENCES `skills`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `user_skills` ADD CONSTRAINT `fk_user_skills_user_id` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE;
