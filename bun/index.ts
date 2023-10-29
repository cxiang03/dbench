import figlet from "figlet";
import { Sequelize, QueryTypes } from "sequelize";

const sequelize = new Sequelize("mysql://root:password@127.0.0.1:3306/dbench");

try {
  await sequelize.authenticate();
  console.log("connection has been established successfully.");
} catch (error) {
  console.error("unable to connect to the database:", error);
}

const server = Bun.serve({
  async fetch(req) {
    const url = new URL(req.url);
    switch (url.pathname) {
      case "/":
        const rst = await sequelize.query("SELECT NOW();", { type: QueryTypes.SELECT });
        if (rst.length === 0) {
          return new Response();
        }
        const now = rst[0];
        const body = figlet.textSync("bun!");
        return new Response(body + "\n" + JSON.stringify(now, null, 2));
      default:
        return new Response();
    }
  },
  port: 3000,
});

console.log(`listening on http://localhost:${server.port}`);
