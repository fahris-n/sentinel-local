local max_tokens = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])

local now = tonumber(redis.call("TIME")[1])
local num_tokens = redis.call("HGET", KEYS[1], "tokens")
local last_accessed = redis.call("HGET", KEYS[1], "last_accessed")

if num_tokens == false or last_accessed == false then
  num_tokens = max_tokens
  last_accessed = now
else
  num_tokens = tonumber(num_tokens)
  last_accessed = tonumber(last_accessed)
end

local elapsed_time = now - last_accessed
local token_refill_amount = elapsed_time * refill_rate
num_tokens = math.min(max_tokens, num_tokens + token_refill_amount)

if num_tokens >= 1 then
  num_tokens = num_tokens - 1
  last_accessed = now
  redis.call("HSET", KEYS[1], "tokens", num_tokens)
  redis.call("HSET", KEYS[1], "last_accessed", last_accessed)
  return 1
else
  redis.call("HSET", KEYS[1], "tokens", num_tokens)
  redis.call("HSET", KEYS[1], "last_accessed", last_accessed)
  return 0
end
