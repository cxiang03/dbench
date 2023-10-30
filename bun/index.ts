import figlet from "figlet";
import { exit } from "process";
import { Sequelize, QueryTypes } from "sequelize";
import { now, randomByPrice } from "./handler";

const hi = figlet.textSync("hi, this is bun!");
console.log(hi);

const dsn = Bun.env.DSN ?? "mysql://root:password@127.0.0.1:3306/dbench";
const sequelize = new Sequelize(dsn);
try {
    await sequelize.authenticate();
    console.log("connection has been established successfully");
} catch (error) {
    console.error("unable to connect to the database:", error);
    exit(1);
}

const server = Bun.serve({
    async fetch(req) {
        const url = new URL(req.url);
        switch (url.pathname) {
            case "/":
                return await now(sequelize);
            case "/random-by-price":
                return await randomByPrice(sequelize);
            default:
                return new Response();
        }
    },
    port: 3000,
});

console.log(`listening on http://localhost:${server.port}`);
