/**
 *  @description
 *      爬取数据
 * 
 */

const requests = require("../libs/requests");



const real_time_task = () => {

    const wait_codes = [ "6005791", "0025452" ];
    const handle_codes = [];

    const timer = setInterval(async () => {

        while (wait_codes.length){

            const code = wait_codes.shift();

            requests("http://mdfm.eastmoney.com", "/EM_UBG_MinuteApi/Js/Get", {
                dtype: 25,
                style: "tail",
                check: "st",
                dtformat: "HH:mm:ss",
                cb: "jQuery1830005435708003650452_1528091676115",
                id: code,
                num: 1,
                _: Date.now()
            }).then(data => {

                console.log("请求回来的数据");
                console.log(data);

            }).catch(err => {

                console.log("发生了错误 ...");
                console.log(err);

            });
            
            await new Promise(resolve => {
                setTimeout(resolve, 1000);
            });
        }

    }, 10000);

};
real_time_task.description = "实时股价 ...";




module.exports = {
    real_time_task
};
