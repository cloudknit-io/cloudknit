module.exports = {
  "type": "mysql",
  "host": process.env.TYPEORM_HOST || "mysqldb",
  "port": 3306,
  "username": process.env.TYPEORM_USERNAME || "root",
  "password": process.env.TYPEORM_PASSWORD || "password",
  "database": process.env.TYPEORM_DATABASE || "development",
  "entities": ["dist/**/**.entity{.ts,.js}"],
  "synchronize": true
};