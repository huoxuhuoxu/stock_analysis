/**
 *  @description
 *      主入口
 * 
 */

require("dotenv").config();

const assert = require("assert");

const app = require("./app/main");
const { test } = require("./libs/tools");
const { info } = require("./libs/logs");

const PORT = ~~process.env.PORT;


const listening = () => {
    info(`已启动, 监听 PORT: ${PORT} ...`);
};


{
    assert(test.isNumber(PORT) && PORT > 1000, "监听的端口号无效");

    app.listen(PORT, listening);
    
}

const { real_time_task } = require("./spider/main");

{
    info(real_time_task.description);
    real_time_task();
}


