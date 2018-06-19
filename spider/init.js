/**
 * @description
 *      初始化, 爬取 股票编号 - 东方财富网股票id - 股票名称
 * 
 */

const fs = require("fs");
const path = require("path");

const requests = require("nodejs-requests");
const { info } = require("../libs/logs");

{
    const start_times = Date.now();
    
    const page_index = 1;
    const page_count = 10000;
    const init_url = `http://nufm.dfcfw.com/EM_Finance2014NumericApplication/JS.aspx?cb=jQuery112407378636087109598_1528454080519&type=CT&token=4f1862fc3b5e77c150a2b985b12db0fd&js=(%7Bdata%3A%5B(x)%5D%2CrecordsTotal%3A(tot)%2CrecordsFiltered%3A(tot)%7D)&cmd=C._A&sty=FCOIATC&st=(ChangePercent)&sr=-1&p=${page_index}&ps=${page_count}&_=1528454080520`;

    (async () => {

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

    })();

}

