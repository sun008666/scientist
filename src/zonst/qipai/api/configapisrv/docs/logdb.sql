-- 日志相关数据表结构，位于logdb

-- 开启关闭游戏平台日志
-- DROP TABLE IF EXISTS game_forbid_set_log;
CREATE TABLE game_forbid_set_log(
    id BIGSERIAL PRIMARY KEY,
    game_id INT NOT NULL DEFAULT 0,
    user_id INT NOT NULL DEFAULT 0,
    username VARCHAR(100) NOT NULL DEFAULT '',
    before_forbid INT NOT NULL DEFAULT 0,
    after_forbid INT NOT NULL DEFAULT 0,
    create_date DATE DEFAULT ('now'::text)::date NOT NULL,
    create_time INTEGER DEFAULT date_part('epoch'::text, now()) NOT NULL
);
COMMENT ON COLUMN game_forbid_set_log.game_id IS '游戏ID';
COMMENT ON COLUMN game_forbid_set_log.user_id IS '用户ID，一般是工号';
COMMENT ON COLUMN game_forbid_set_log.username IS '用户名';
COMMENT ON COLUMN game_forbid_set_log.before_forbid IS '修改之前的状态，1表示已禁用，0表示未禁用';
COMMENT ON COLUMN game_forbid_set_log.after_forbid IS '修改之后的状态，1表示已禁用，0表示未禁用';
COMMENT ON COLUMN game_forbid_set_log.create_date IS '创建日期';
COMMENT ON COLUMN game_forbid_set_log.create_time IS '创建时间';
CREATE INDEX "game_forbid_set_log_game_idx" ON "game_forbid_set_log" (game_id);
CREATE INDEX "game_forbid_set_log_user_id_idx" ON "game_forbid_set_log" (user_id);

