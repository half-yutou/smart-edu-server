
-- 插入基础数据：学科
INSERT INTO subjects (name) VALUES 
('数学'), ('语文'), ('英语'), ('物理'), ('化学'), ('生物'), ('历史'), ('地理'), ('政治');

-- 插入基础数据：年级
INSERT INTO grades (name) VALUES 
('一年级'), ('二年级'), ('三年级'), ('四年级'), ('五年级'), ('六年级'), 
('初一'), ('初二'), ('初三'), 
('高一'), ('高二'), ('高三');

-- 插入测试用户
-- 密码均为 123456 (加密后示例，实际应使用 bcrypt 生成)
-- 假设 bcrypt(123456) = $2a$10$x.z.y... (这里为了演示简便，假设应用层会处理加密，或者手动生成一个合法的hash)
-- 这里使用一个真实的 bcrypt hash 值: $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy (对应 123456)

INSERT INTO users (username, password, role, nickname, avatar_url) VALUES 
('teacher1', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'teacher', '张老师', 'https://example.com/avatar/t1.jpg'),
('teacher2', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'teacher', '李老师', 'https://example.com/avatar/t2.jpg'),
('student1', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'student', '小明', 'https://example.com/avatar/s1.jpg'),
('student2', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'student', '小红', 'https://example.com/avatar/s2.jpg'),
('student3', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'student', '小刚', 'https://example.com/avatar/s3.jpg');

-- 插入测试班级
-- 注意：teacher_id, subject_id, grade_id 需要对应上面的 ID
-- 假设 users 表 id 从 1 开始: teacher1=1, teacher2=2
-- 假设 subjects 表: 数学=1, 语文=2
-- 假设 grades 表: 初一=7, 高一=10

INSERT INTO classes (name, code, teacher_id, subject_id, grade_id, description) VALUES 
('初一(1)班 数学提高班', 'ABC123', 1, 1, 7, '数学竞赛冲刺'),
('高一(2)班 语文兴趣组', 'XYZ789', 2, 2, 10, '古诗词鉴赏');

-- 插入班级成员
-- student1(3) 加入 class 1
-- student2(4) 加入 class 1
-- student3(5) 加入 class 2
INSERT INTO class_members (class_id, student_id) VALUES 
(1, 3),
(1, 4),
(2, 5);

-- 插入资源
INSERT INTO resources (title, description, res_type, file_url, subject_id, grade_id, uploader_id, duration) VALUES 
('函数的基本概念', '讲解函数的定义域和值域', 'video', 'https://oss.example.com/video/math_func.mp4', 1, 10, 1, 1200),
('唐诗三百首赏析', '精选唐诗讲解', 'document', 'https://oss.example.com/doc/tang_poem.pdf', 2, 7, 2, 0);

-- 班级关联资源
INSERT INTO class_resources (class_id, resource_id) VALUES 
(1, 1), -- 班级1 关联 资源1
(2, 2); -- 班级2 关联 资源2

-- 插入作业
INSERT INTO homeworks (title, class_id, creator_id, deadline) VALUES 
('第一章函数作业', 1, 1, '2025-12-31 23:59:59+08');

-- 插入题目
INSERT INTO questions (homework_id, question_type, content, options, correct_answer, score, order_num) VALUES 
(1, 'choice', '下列哪个函数是偶函数？', '{"A": "y=x", "B": "y=x^2", "C": "y=x^3", "D": "y=sin(x)"}', 'B', 5, 1),
(1, 'text', '请简述函数的单调性定义。', NULL, '略', 10, 2);

-- 插入提交记录
INSERT INTO submissions (homework_id, student_id, status, total_score) VALUES 
(1, 3, 'graded', 15);

-- 插入答题详情
INSERT INTO submission_details (submission_id, question_id, answer_content, is_correct, score) VALUES 
(1, 1, 'B', true, 5),
(1, 2, '当x1<x2时...', true, 10);
