
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
  return db.addColumn("shares", "sign", {
    type: "smallint",
    length: 2,
    notNull: true
  }, callback);
};

exports.down = function(db) {
  return db.removeColumn("shares", "sign");
};

exports._meta = {
  "version": 1
};
