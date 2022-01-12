
local ttl_window = tonumber(ARGV[1])
redis.debug(KEYS[1])
redis.debug(ttl_window)
if redis.call("EXISTS", KEYS[1]) == 1 then
	local eid = redis.call("hget", KEYS[1], "eid")
	redis.debug(eid)
	return eid
else
	return ""
end
