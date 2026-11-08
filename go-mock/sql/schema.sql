-- go-mock 建表语句（幂等，可重复执行）

-- 用户表
CREATE TABLE IF NOT EXISTS user (
  uid    TEXT PRIMARY KEY,            -- 用户 id
  uname  TEXT NOT NULL,               -- 昵称
  avatar TEXT NOT NULL DEFAULT ''     -- 头像地址
);

-- 评论表
CREATE TABLE IF NOT EXISTS comment (
  rpid       INTEGER PRIMARY KEY AUTOINCREMENT, -- 评论 id
  uid        TEXT    NOT NULL REFERENCES user(uid),
  content    TEXT    NOT NULL,                  -- 评论内容
  like_count INTEGER NOT NULL DEFAULT 0,        -- 点赞数
  ctime      INTEGER NOT NULL,                  -- 发布时间(unix 秒)
  is_deleted INTEGER NOT NULL DEFAULT 0         -- 软删除标记 0/1
);

CREATE INDEX IF NOT EXISTS idx_comment_ctime ON comment(ctime DESC);
CREATE INDEX IF NOT EXISTS idx_comment_like  ON comment(like_count DESC);
