/**
 *  @readme
 *      连接mysql
 *      
 *      QUERY: 通用
 *      TRANSACTION: 事务
 * 
 */

require("dotenv").config();

const mysql = require('mysql');
const mysql_config = {
    connectionLimit: 10,
    host: process.env.DB_HOST,
    post: process.env.DB_PORT,
    user: process.env.DB_USERNAME,
    password: process.env.DB_PASSWORD,
    database: process.env.DB_DATABASE,
    multipleStatements: true
};
let pool = mysql.createPool(mysql_config);

// 通用语句
const QUERY = function (...argv){
    return new Promise(function(resolve, reject){
        pool.getConnection((err, conn) => {
            if(err){
                err.status = 12001;
                throw err;
            }
            conn.query(...argv, function(err, rows, fields){
                conn.release();
                if(err){
                    err.status = 12002;
                    reject(err);
                    return ;
                }
                resolve(rows);
            });
        });
    });
};

// 事务
const TRANSACTION = async () => {
  
    let connect = await new Promise(resolve => {
        pool.getConnection((err, conn) => {
            if (err) {
                err.status = 12101;
                throw err;
            }
            resolve(conn);
        }); 
    });

    return {
        next: (...argv) => {
            if (!argv.length) {
                connect.release();
                let ret = {
                    value: undefined,
                    done: true
                };
                return ret;
            }
            return new Promise((resolve, reject) => {
                connect.query(...argv, (err, rows, fields) => {
                    if (err) {
                        return reject(err);
                    }
                    let ret = {
                        value: rows,
                        done: false
                    };
                    resolve(ret);
                });
            });
        }
    };

};



module.exports = {
    QUERY,
    TRANSACTION
};