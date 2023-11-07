import figlet from "figlet";
import { exit } from "process";
import { Sequelize } from "sequelize";
import { now, es, randomByPrice, randomInTimeRange, esInTimeRange, randomUUID, meiliInTimeRange } from "./handler";
import { Client } from "es7";
import { createClient, RedisClientType } from "redis";
import { MeiliSearch } from "meilisearch";

const hi = figlet.textSync("hi, this is bun!");
console.log(hi);

const dsn = Bun.env.DSN ?? "mysql://root:password@127.0.0.1:3306/playground";
const sequelize = new Sequelize(dsn);
try {
    await sequelize.authenticate();
    console.log("connection has been established successfully");
} catch (error) {
    console.error("unable to connect to the database:", error);
    exit(1);
}

const redis: RedisClientType = createClient();
redis.on("error", (err) => console.log("redis client error", err)).connect();

// const meili = new MeiliSearch({
//     host: "http://103.3.60.74:7700",
//     apiKey: "123qQ123",
// });

// const esa = Bun.env.ESA ?? "http://127.0.0.1:9200";
// const esClient = new Client({ node: esa });
// try {
//     await esClient.ping();
//     console.log("connection has been established successfully");
// } catch (error) {
//     console.error("unable to connect to the elastic_search:", error);
//     exit(1);
// }

const server = Bun.serve({
    async fetch(req) {
        const url = new URL(req.url);
        switch (url.pathname) {
            case "/":
                return await now(sequelize);
            // case "/es":
            //     return await es(esClient);
            // case "/es-in-time-range":
            //     return await esInTimeRange(esClient);
            case "/random-by-price":
                return await randomByPrice(sequelize);
            case "/random-in-time-range":
                return await randomInTimeRange(sequelize);
            case "/random-uuid":
                return await randomUUID(redis);
            // case "/meili-in-time-range":
            //     return await meiliInTimeRange(meili);
            default:
                return new Response(JSON.stringify({ status: "ok" }, null, 2));
        }
    },
    port: 3000,
});

console.log(`listening on http://localhost:${server.port}`);
