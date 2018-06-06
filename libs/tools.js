/**
 *  @description
 *      提供基础方法
 *  
 *  @function
 *      test.isNumber: NaN视为非数字类型
 *           
 * 
 */

const test = (() => {

    const type_arr = [
        "Function",
        "Object",
        "Array",
        "Symbol",
        "Set",
        "Map",
        "WeakSet",
        "WeakMap",
        "String",
        "Number",
        "Boolean"
    ];

    const tmp = {};

    for ( const type_name of type_arr){
        tmp[`is${type_name}`] = (v) => {
            return Object.prototype.toString.call(v) === `[object ${type_name}]`;
        };
    }

    // isNumber 方法 - 追加条件
    const tmp_number = tmp["isNumber"];
    tmp["isNumber"] = (v) => tmp_number(v) && !Number.isNaN(v);

    return tmp;

})();


module.exports = {
    test
};

