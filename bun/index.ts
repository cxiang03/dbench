import figlet from "figlet";
import { exit } from "process";
import { Sequelize } from "sequelize";
import { now, es, randomByPrice, randomInTimeRange, esInTimeRange } from "./handler";
import { Client } from "es7";

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

const esa = Bun.env.ESA ?? "http://127.0.0.1:9200";
const esClient = new Client({ node: esa });
try {
    await esClient.ping();
    console.log("connection has been established successfully");
} catch (error) {
    console.error("unable to connect to the elastic_search:", error);
    exit(1);
}

const server = Bun.serve({
    async fetch(req) {
        const url = new URL(req.url);
        switch (url.pathname) {
            case "/":
                return await now(sequelize);
            case "/es":
                return await es(esClient);
            case "/random-by-price":
                return await randomByPrice(sequelize);
            case "/random-in-time-range":
                return await randomInTimeRange(sequelize);
            case "/es-in-time-range":
                return await esInTimeRange(esClient);
            default:
                return new Response();
        }
    },
    port: 3000,
});

console.log(`listening on http://localhost:${server.port}`);
