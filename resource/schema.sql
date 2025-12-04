-- SmartEduHub Database Schema
-- Database: smarteduhub
-- Note: 物理外键约束已移除，由应用层保证数据一致性。

-- ==========================================
-- 1. 用户中心 (User Center)
-- ==========================================

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,  -- 用户名 (纯数字唯一)
    password VARCHAR(255) NOT NULL, -- 加密后的密码
    role VARCHAR(20) NOT NULL,      -- student, teacher, admin
    nickname VARCHAR(50) NOT NULL DEFAULT '',
    avatar_url VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
CREATE UNIQUE INDEX idx_users_username ON users(username);

-- 学科表
CREATE TABLE IF NOT EXISTS subjects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
    );
CREATE UNIQUE INDEX idx_subjects_name ON subjects(name);

-- 年级表
CREATE TABLE IF NOT EXISTS grades (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
    );
CREATE UNIQUE INDEX idx_grades_name ON grades(name);


-- ==========================================
-- 2. 班级管理 (Class Management)
-- ==========================================

-- 班级表
CREATE TABLE IF NOT EXISTS classes (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(10) NOT NULL,      -- 6位班级邀请码
    teacher_id BIGINT NOT NULL,     -- 关联 users.id
    subject_id INT DEFAULT 0,       -- 关联 subjects.id (可选)
    grade_id INT DEFAULT 0,         -- 关联 grades.id
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
CREATE UNIQUE INDEX idx_classes_code ON classes(code);
CREATE INDEX idx_classes_teacher_id ON classes(teacher_id);

-- 班级成员表 (学生-班级关联)
CREATE TABLE IF NOT EXISTS class_members (
    id BIGSERIAL PRIMARY KEY,
    class_id BIGINT NOT NULL,       -- 关联 classes.id
    student_id BIGINT NOT NULL,     -- 关联 users.id
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX idx_class_members_unique ON class_members(class_id, student_id);
CREATE INDEX idx_class_members_student_id ON class_members(student_id);


-- ==========================================
-- 3. 资源管理 (Resource Management)
-- ==========================================

-- 资源表
CREATE TABLE IF NOT EXISTS resources (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    res_type VARCHAR(20) NOT NULL,  -- video, document
    file_url VARCHAR(255) NOT NULL, -- OSS URL
    subject_id INT NOT NULL,        -- 关联 subjects.id
    grade_id INT NOT NULL,          -- 关联 grades.id
    uploader_id BIGINT NOT NULL,    -- 关联 users.id
    duration INT DEFAULT 0,         -- 视频时长(秒)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
CREATE INDEX idx_resources_subject_grade ON resources(subject_id, grade_id);
CREATE INDEX idx_resources_uploader_id ON resources(uploader_id);

-- 班级资源关联表
CREATE TABLE IF NOT EXISTS class_resources (
    id BIGSERIAL PRIMARY KEY,
    class_id BIGINT NOT NULL,       -- 关联 classes.id
    resource_id BIGINT NOT NULL,    -- 关联 resources.id
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX idx_class_resources_unique ON class_resources(class_id, resource_id);

-- 弹幕表
CREATE TABLE IF NOT EXISTS danmakus (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL,    -- 关联 resources.id
    user_id BIGINT NOT NULL,        -- 关联 users.id
    content VARCHAR(255) NOT NULL,
    time_point FLOAT NOT NULL,      -- 视频时间点(秒)
    color VARCHAR(10) DEFAULT '#FFFFFF',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
CREATE INDEX idx_danmakus_resource_time ON danmakus(resource_id, time_point);


-- ==========================================
-- 4. 作业系统 (Homework System)
-- ==========================================

-- 作业表
CREATE TABLE IF NOT EXISTS homeworks (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    class_id BIGINT NOT NULL,       -- 关联 classes.id
    creator_id BIGINT NOT NULL,     -- 关联 users.id
    deadline TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
CREATE INDEX idx_homeworks_class_id ON homeworks(class_id);

-- 题目表
CREATE TABLE IF NOT EXISTS questions (
    id BIGSERIAL PRIMARY KEY,
    homework_id BIGINT NOT NULL,        -- 关联 homeworks.id
    question_type VARCHAR(20) NOT NULL, -- choice(选择), text(主观)
    content TEXT NOT NULL,
    options JSONB,                      -- 选项: {"A":"..", "B":".."}
    correct_answer TEXT,                -- 参考答案
    score INT NOT NULL,                 -- 分值
    order_num INT DEFAULT 0,            -- 排序
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
CREATE INDEX idx_questions_homework_id ON questions(homework_id);

-- 提交记录表 (一次作业提交)
CREATE TABLE IF NOT EXISTS submissions (
    id BIGSERIAL PRIMARY KEY,
    homework_id BIGINT NOT NULL,    -- 关联 homeworks.id
    student_id BIGINT NOT NULL,     -- 关联 users.id
    status VARCHAR(20) DEFAULT 'submitted', -- submitted, graded
    total_score INT DEFAULT 0,
    feedback TEXT NOT NULL DEFAULT '',                  -- 教师总评
    ai_feedback TEXT NOT NULL DEFAULT '',               -- 教师AI评语
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
CREATE UNIQUE INDEX idx_submissions_unique ON submissions(homework_id, student_id);

-- 答题详情表 (每一道题的回答)
CREATE TABLE IF NOT EXISTS submission_details (
    id BIGSERIAL PRIMARY KEY,
    submission_id BIGINT NOT NULL,  -- 关联 submissions.id
    question_id BIGINT NOT NULL,    -- 关联 questions.id
    answer_content TEXT NOT NULL DEFAULT '',            -- 学生答案 (OCR后或手动)
    answer_image_url VARCHAR(255) NOT NULL DEFAULT '',  -- 答案原图
    is_correct BOOLEAN DEFAULT FALSE,
    score INT DEFAULT 0,
    comment VARCHAR(255) NOT NULL DEFAULT ''            -- 教师单题评语
    );
CREATE INDEX idx_submission_details_submission_id ON submission_details(submission_id);