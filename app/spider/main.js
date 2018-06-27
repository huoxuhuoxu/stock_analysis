/**
 *  @description
 *      爬取数据
 * 
 */

const requests = require("nodejs-requests");
const { info } = require("../libs/logs");

/**
 *  @TODO
 *      换方案, 这种形式, 会产生 一定时间内没有处理完, 第二轮开始的问题 ...
 *      解决方案 requests - 方法内部增加 超时字段, setInterval 换 setTimeout
 *      requests 内部使用的promise需要 增加finally 方法 ...
 */

// const wait_codes = [ "6005791", "0025452" ];
const wait_codes = [];
const handle_codes = new Set();

const stock_information = async () => {
    while (wait_codes.length){

        const code = wait_codes.shift();

        requests("http://mdfm.eastmoney.com/EM_UBG_MinuteApi/Js/Get", {
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

            handle_codes.add(code);

        }).catch(err => {

            console.log("发生了错误 ...");
            console.log(err);

            handle_codes.add(code);

        });
        
        await new Promise(resolve => {
            setTimeout(resolve, 1000);
        });
    }

    wait_codes.push(...Array.from(handle_codes));
    handle_codes.clear();
};


const A_shares = async () => {

    const data = await requests("http://nufm.dfcfw.com/EM_Finance2014NumericApplication/JS.aspx?type=CT&cmd=P.(x),(x)|0000011|3990012&sty=SHSTD|SZSTD&st=z&token=4f1862fc3b5e77c150a2b985b12db0fd&cb=jQuery112407378636087109598_1528454080515&_=1528454080516")

    console.log("大盘数据:");
    const [ A_hu, A_shen ] = JSON.parse(data.match(/\[.*\]/)[0]).map(v => v.split(","));
    const info = [
        {
            title: "港股通(沪)",
            list: A_hu.slice(0, 6)
        },
        {
            title: "沪股通",
            list: A_hu.slice(6)
        },
        {
            title: "港股通(深)",
            list: A_shen.slice(0, 6)
        },
        {
            title: "深股通",
            list: A_shen.slice(6)
        }
    ];

    let diff = 0;
    for (const item of info){
        console.log(`${item.title} 流入 ${item.list[2]} 流出 ${item.list[1]} 差值 ${item.list[0]}`);
        diff += item.list[0];
    }
    console.log("\r\n结语:");
    console.log(diff > 0 ? "外盘资金在进场" : "外盘资金在逃");

};


const real_time_task = () => {

    const timer = setInterval(async () => {

        stock_information();
        A_shares();

    }, 10000);

    process.nextTick(() => {
        info("启动爬虫中, 请稍候 ...");
    });

};
real_time_task.description = "实时股价 ...";





module.exports = {
    real_time_task
};
