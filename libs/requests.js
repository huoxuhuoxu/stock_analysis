/** 
 * @description 
 *      提供 请求方法
 * 
 */

// requests
const http = require("http");
const https = require("https");
const { URL } = require("url");
const assert = require("assert");


/**
 * @description
 *      get 方式发起请求
 * 
 * @param {*} protocol          协议, .eg: http, https
 * @param {*} url               完整的请求地址
 * @param {*} callback          成功后的回调函数
 * 
 * @return Promise
 */
const get = async (protocol, url, callback) => {
    return new Promise((resolve, reject) => {
        protocol.get(url, callback.bind(null, resolve, reject))
                .on("error", (e) => {
                    reject(e);
                });
    });
};


/**
 * @description 
 *      post 方式发起请求
 * 
 * @param {*} protocol 
 * @param {*} options           请求头信息
 * @param {*} data              需要发送的数据（请求体）
 * @param {*} callback          
 * 
 * @return Promise
 */
const post = async (protocol, options, data, callback) => {
    return new Promise((resolve, reject) => {
        const req = protocol.request(options, callback.bind(null, resolve, reject));
        req.on("error", (e) => { reject(e); });
        req.write(data);
        req.end();
    }); 
};


/**
 * @description
 *      根据参数生成 请求头与请求体     
 * 
 * @param {*} origin                域名 
 * @param {*} path                  具体路径
 * @param {*} body                  数据
 * @param {*} method                发请请求的方式, .eg: post, put, delete, option, ...
 * @param {*} req_content_type      发送数据时的数据格式化方式
 * 
 * @return [ 请求头, 请求体 ]
 */
const getOptions = (origin, path, body, method, req_content_type) => {
    const originArr = origin.split("://");
    const pathArr = originArr[1].split(":");
    let hostname = pathArr[0], 
        port = originArr[0] === "http" ? 80 : 443;
    if (pathArr.length > 1) {
        port = pathArr[1];
    }

    const postData = JSON.stringify(body);
    const options = {
        hostname,
        port,
        path,
        method,
        headers: {
            "Content-Type": req_content_type,
            "Content-Length": Buffer.byteLength(postData)
        }
    };
    return [ options,  postData];
};


/**
 * @description
 *      处理返回的数据
 * 
 * @param {*} content_type      验证接受的数据采用的格式化方式, .eg: json, txt, ...
 * @param {*} resolve           Promise.resolve
 * @param {*} reject            Promise.reject
 * @param {*} res               response object
 * 
 * @return undefined
 */
const callback = async (content_type, resolve, reject, res) => {

    const { statusCode } = res;
    const contentType = res.headers["content-type"];

    let error;
    if (statusCode !== 200){
        error = new Error("请求失败\n" + `状态码: ${statusCode}`);
    } else if ( ! (new RegExp(content_type).test(contentType)) ){
        error = new Error("无效 content-type\n" + `期望: ${content_type}, 实际获取: ${contentType}`);
    }

    if (error){
        res.resume();
        reject(error);
        return ;
    }

    const chunks = [];
    res.setEncoding("utf8");
    res.on("data", chunk => {
        if (chunk !== null){
            chunks.push(chunk);
        }
    });
    res.on("end", () => {
        resolve(chunks.join(""));
    });
    res.on("error", (e) => {
        console.error("[error] 错误: ", error.message);
        reject(e);
    });
};



/**
 * @description
 *      同一处理 协议, 发起请求的方法等
 * 
 * @param {*} origin        域名
 * @param {*} pathname      具体路径
 * @param {*} body          数据
 * @param {*} method        方法
 * @param {*} param4        
 *      res_content_type: 返回时需要验证返回数据的格式化方式, 不需要验证 ''
 *      req_content_type: 发起请求时存在请求体情况下, 数据的格式化方式
 * 
 * @return Promise
 */
module.exports = async (origin, pathname, body = {}, method = "GET", 
    { 
        res_content_type = "",
        req_content_type = "application/json"
    } = { 
        res_content_type: "",
        req_content_type: "application/json"
    }
) => {

    assert(/https?:\/\/.*/.test(origin), "origin-format error");

    let protocol = http;
    /^https/.test(origin) && (protocol = https);

    const cb = callback.bind(null, res_content_type);

    if (method.toLowerCase() === "get") {
        let search = "";
        for (let [k, v] of Object.entries(body)){
            search += `&${k}=${v}`;
        }
        return await get(protocol, `${new URL(pathname, origin).href}?${search.substr(1)}`, cb);
    }

    const [ options, data ] = getOptions(origin, pathname, body, method, req_content_type);
    return await post(protocol, options, data, cb);

}; 




