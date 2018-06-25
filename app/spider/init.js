/**
 * @description
 *      初始化, 爬取 股票编号 - 东方财富网股票id - 股票名称
 * 
 */

const fs = require("fs");
const path = require("path");
const assert = require("assert");

const argv = require("optimist").default({
    "s": "init"
}).argv;
const requests = require("nodejs-requests");

const { QUERY } = require("../db/connect_mysql");
const { logs: { info, warn }, tools: { test } } = require("../libs");


const QUERYS = async (sql_str, sql_list, num = 100) => {
    assert(
        test.isArray(sql_list) && 
        test.isString(sql_str) &&
        test.isNumber(num), 
        "exec sql, but params error !"
    );

    let i=0, j, exec_list;
    do {
        j = i + 1;
        exec_list = sql_list.slice(i * num, j * num);
        const exec_str = Array(exec_list.length).fill(sql_str).join(";");

        const list = [];
        for (const v of exec_list){
            list.push(...v);
        }
        
        await QUERY(exec_str, list);

        i = j;

    } while (exec_list.length === num);

};

const init = async (json) => {
    const sql_list = [];
    for (let [ key, value ] of Object.entries(json)){
        sql_list.push([ value.name, key, value.sign ]);
    }
    await QUERYS(
        "insert into shares (name, number, sign) values (?, ?, ?)",
        sql_list
    );
};

const update = async (json) => {
    const ret = await QUERY("select number from shares");
    const numbers = new Set();
    for (const v of ret){
        numbers.add(v["number"]);
    }

    for (let [ key, value ] of Object.entries(json)){
        if (!numbers.has(key)){
            const obj = {};
            obj[key] = value;
            await init(obj);
            numbers.add(key);
            info("新增: ", key);
        }
    }
};


{

    const start_times = Date.now();
    
    const page_index = 1;
    const page_count = 10000;
    const init_url = `http://nufm.dfcfw.com/EM_Finance2014NumericApplication/JS.aspx?cb=jQuery112407378636087109598_1528454080519&type=CT&token=4f1862fc3b5e77c150a2b985b12db0fd&js=(%7Bdata%3A%5B(x)%5D%2CrecordsTotal%3A(tot)%2CrecordsFiltered%3A(tot)%7D)&cmd=C._A&sty=FCOIATC&st=(ChangePercent)&sr=-1&p=${page_index}&ps=${page_count}&_=1528454080520`;

    (async () => {

        if ( module.parent ){
            module.exports = {};
            return ;
        }

        const s_txt = await requests(init_url);
        const resul = (s_txt.split("data:[\"")[1]).split("\"],");

        const data = resul[0];
        const other = resul[1];

        const datas = data.split("\",\"");
        const json = {};

        for (const s of datas){
            const arr = s.split(",");
            json[arr[1]] = {
                name: arr[2],
                sign: arr[0]
            };
        }

        const [ , total ] = other.match(/recordsTotal:(\d+),/);
        json['total'] = total;

        fs.writeFileSync(path.resolve(__dirname, "data.json"), JSON.stringify(json, null, "\t"));

        console.log("\r\n");
        info("初始化A股信息成功, 目前A股一共 %d 只股票\r\n", total);
        info("耗时: %d ms", Date.now() - start_times );
        console.log("\r\n");

        info("开始写入数据库 \r\n");
        delete json['total'];

        switch (argv.s){
            case "init":
                await init(json);
            break;
            case "update":
                await update(json);
            break;
            default:
                warn("未知参数: ", argv.s);
        }
        info("end ... \r\n");

        process.exit(0);

    })();

}

