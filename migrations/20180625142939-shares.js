
var dbm;
var type;
var seed;

/**
  * We receive the dbmigrate dependency from dbmigrate initially.
  * This enables us to not have to rely on NODE_PATH.
  */
exports.setup = function(options, seedLink) {
  dbm = options.dbmigrate;
  type = dbm.dataType;
  seed = seedLink;
  console.log(dbm, type, seed);
};

exports.up = function(db, callback) {
  return db.createTable("shares", {
    id: {
      type: "int",
      primaryKey: true,
      autoIncrement: true,
      notNull: true
    },
    name: {
      type: "string",
      notNull: true
    },
    number: {
      type: "string",
      notNull: true,
      length: 8
    },
    updated_at: "datetime",
    created_at: "datetime"
  }, function (){
    db.addIndex("shares", "shares_name", "name", true);
    db.addIndex("shares", "shares_number", "number", true);
    callback();
  });

};

exports.down = function(db) {
  return db.dropTable("shares");
};

exports._meta = {
  "version": 1
};
