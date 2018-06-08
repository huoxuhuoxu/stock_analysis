/** 
 * @description 
 *      提供 请求方法
 * 
 */

// requests
const http = require("http");
const https = require("https");
const url = require("url");
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
 * @param {*} url_info              URL对象
 * @param {*} body                  数据
 * @param {*} method                发请请求的方式, .eg: post, put, delete, option, ...
 * @param {*} req_content_type      发送数据时的数据格式化方式
 * 
 * @return [ 请求头, 请求体 ]
 */
const getOptions = (url_info, body, method, req_content_type) => {
    const postData = JSON.stringify(body);
    const options = {
        hostname: url_info.hostname,
        port: url_info.port,
        path: url.path,
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
 * @param {*} real_url      完整的URL
 * @param {*} body          数据
 * @param {*} method        方法
 * @param {*} param4        
 *      res_content_type: 返回时需要验证返回数据的格式化方式, 不需要验证 ''
 *      req_content_type: 发起请求时存在请求体情况下, 数据的格式化方式
 * 
 * @desc
 *      调用时只传了real_url, 直接执行 http(s).get(real_url)
 * 
 * 
 * @return Promise
 */
module.exports = async (real_url, body = {}, method = "GET", 
    { 
        res_content_type = "",
        req_content_type = "application/json"
    } = { 
        res_content_type: "",
        req_content_type: "application/json"
    }
) => {

    assert(real_url, "request arguments");

    const url_info = url.parse(real_url);
    assert(/https?:/.test(url_info.protocol), "origin-format error, can only be http, https");

    let protocol = http;
    /^https/.test(url_info.protocol) && (protocol = https);

    const cb = callback.bind(null, res_content_type);

    if (method.toLowerCase() === "get") {

        if (JSON.stringify(body) === "{}"){
            return await get(protocol, real_url, cb);
        } 

        let search = "";
        for (let [k, v] of Object.entries(body)){
            search += `&${k}=${v}`;
        }
        return await get(protocol, `${url_info.protocol}//${url_info.host}${url_info.pathname}?${search.substr(1)}`, cb);
    
    }

    const [ options, data ] = getOptions(url_info, body, method, req_content_type);
    return await post(protocol, options, data, cb);

}; 




