/**
 *  @description
 *      各种类型的输出
 * 
 */

const util = require("util");


const outpout = (...args) => {

    const [ color ] = args.splice(0, 1);
    const s_format = util.format(...args);
    console.log(color, s_format);

};

const log = (...args) => {
    outpout("\x1b[0m", ...args, "\x1b[0m");
};

const info = (...args) => {
    outpout("\x1b[32m", ...args, "\x1b[0m");
};

const warn = (...args) => {
    outpout("\x1b[33m", ...args, "\x1b[0m");
};

const error = (...args) => {
    outpout("\x1b[91m", ...args, "\x1b[0m");
};


module.exports = {
    log,
    info,
    warn,
    error
};

