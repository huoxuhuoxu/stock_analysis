/**
 *  @description
 *      主入口
 * 
 */

require("dotenv").config();

const assert = require("assert");

const { server, spider } = require("./app");
const { tools: { test }, logs: { info } } = require("./app/libs");

const PORT = ~~process.env.PORT;


const listening = () => {
    info(`已启动 Web Service, 监听 PORT: ${PORT} ...`);
};


{
    assert(test.isNumber(PORT) && PORT > 1000, "监听的端口号无效");
    server.listen(PORT, listening);
}

{
    const { real_time_task } = spider;
    info(real_time_task.description);
    real_time_task();
}


