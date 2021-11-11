
-- 子游戏配置增加桌子数量，每桌座位数量，每桌最小用户人数
ALTER TABLE game_area ADD COLUMN desk_nums integer NOT NULL DEFAULT 4;
ALTER TABLE game_area ADD COLUMN seat_nums integer NOT NULL DEFAULT 4;
ALTER TABLE game_area ADD COLUMN min_user_nums integer NOT NULL DEFAULT 4;


-- 游戏房间服务增加座位数量，最小用户数量
ALTER TABLE game_room_sh ADD COLUMN seat_nums integer DEFAULT 4;
ALTER TABLE game_room_sh ADD COLUMN min_user_nums integer DEFAULT 4;
ALTER TABLE game_room_sh ADD COLUMN master_node_id integer NOT NULL DEFAULT 0;


-- 子游戏服务配置
CREATE TABLE game_area (
    id serial PRIMARY KEY,
    game_id integer NOT NULL,
    game_area_id integer NOT NULL,
    game_area_name text NOT NULL,
    deploy_progame_name text NOT NULL,
    deploy_program_path text NOT NULL DEFAULT '',
    UNIQUE(game_id, game_area_id)
);

-- 房间服务配置
CREATE TABLE game_room_sh (
    id serial PRIMARY KEY,
    name text NOT NULL,
    game_id integer NOT NULL,
    game_area_id integer NOT NULL,
    game_category_id integer NOT NULL DEFAULT 0,
    room_type integer NOT NULL DEFAULT 0,
    room_server_ip varchar(20) NOT NULL,
    room_server_port integer NOT NULL,
    master_node_id integer NOT NULL, -- 主节点ID

    desk_nums integer NOT NULL DEFAULT 250,
    seat_nums integer NOT NULL DEFAULT 4,
    min_user_nums integer NOT NULL DEFAULT 4,
    status integer NOT NULL DEFAULT 1
);
-- status 状态