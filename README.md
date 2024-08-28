# ginRanking


Gin + Gorm + Redis 开发朋友圈投票 点赞 助力


### 提交备忘
    390cd04ca49c8ee8c17ed475a10c91a60864bc90
    完成前后端，页面接口基本已通，排行榜通过Redis未使用，与教程相比有优化和调整 加入了自己的一些设计


### 后续待办
    1. 拆表 score独立单表 涉及所有方法修改  OK
    2. 优化排行榜接口 DB优化 完善Redis排行 结构统一 OK
    3. 增加选手信息Redis缓存 优先读取 分数读取走Zset 
    4. 编码规范 编码优化 变量大小写 驼峰 规则确定 统一
    5. 增加Service层 controller model service (不必须)
    6. Gin 框架高级学习 中间件 拦截器... 优化 框架封装 Redis GORM 封装使用


### 代码参考
    [后端完整](https://github.com/guardian-wjt/gin_ranking)
    [带前端代码](https://github.com/CyberMidori/gin-ranking)

