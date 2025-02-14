local key = KEYS[1]
-- 用户输入的code
local expectedCode = ARGV[1]
local code = redis.call("get",key)
local cntKey = key..":cnt"
local cnt = tonumber( redis.call("get",cntKey))

if cnt == nil or cnt <= 0 then
    -- 说明一直出错，超过验证次数
    -- 或者已经用过了
    return -1
elseif expectedCode == code then
    -- 输入正确
    redis.call("set",cntKey,-1)
    -- redis.call("del",key)
    return 0
else
    redis.call("decr",cntKey)
    return -2
end