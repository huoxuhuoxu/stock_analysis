/**
 *  @description
 *      libs/tools
 * 
 */

const assert = require("assert");
const tools = require("../libs/tools");


describe ("测试 test 提供的方法", () => {

    const { test } = tools;

    it ("通用类型检查", (done) => {

        const type_data = {
            "Function": function (){},
            "Array": [],
            "Number": Math.random() * 1000,
            "String": "",
            "Boolean": true,
            "WeakMap": new WeakMap(),
            "Map": new Map(),
            "WeakSet": new WeakSet(),
            "Set": new Set(),
            "Symbol": Symbol(""),
            "Object": {}
        };

        for (const key in type_data){

            let sign = false;

            for (const type_name in test){

                const b_resul = test[type_name](type_data[key]);
                
                if (key === type_name.substring(2)){
                    assert(b_resul, `类型检查出错, ${type_name}, ${key}`);
                    sign = true;
                    continue;
                }

                assert(!b_resul, `类型检查出错, ${type_name}, ${key}`);

            }

            assert(sign, `未检测出有效类型, ${key}`);

        }

        done();

    });

    it ("NaN类型检查", (down) => {

        const nan = NaN;
        assert(!test.isNumber(nan), "未检测出NaN");
        down();

    });

});


