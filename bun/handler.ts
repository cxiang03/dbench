import { Sequelize, QueryTypes } from "sequelize";
import { Client } from "es7";

async function now(sequelize: Sequelize): Promise<Response> {
    const query = "SELECT NOW()";
    const results = await sequelize.query(query, { type: QueryTypes.SELECT });
    const body = results.shift();
    return new Response(JSON.stringify(body, null, 2));
}

async function es(esClient: Client): Promise<Response> {
    const { body } = await esClient.info();
    return new Response(JSON.stringify(body, null, 2));
}

async function randomByPrice(sequelize: Sequelize): Promise<Response> {
    // magic random numbers
    const from = 1630000000 + Math.floor(Math.random() * 10000000);
    const to = from + Math.floor(Math.random() * 30000000);
    const query = `SELECT * FROM prices WHERE ${from} <= time_stamp AND time_stamp <= ${to} AND p_type ="F" ORDER BY price DESC LIMIT 5`;
    const results = await sequelize.query(query, { type: QueryTypes.SELECT });
    const body = results.shift();
    return new Response(JSON.stringify(body, null, 2));
}

async function randomInTimeRange(sequelize: Sequelize): Promise<Response> {
    // magic random numbers
    const from = 1630000000 + Math.floor(Math.random() * 10000000);
    const to = from + Math.floor(Math.random() * 30000000);
    const query = `SELECT * FROM prices WHERE ${from} <= time_stamp AND time_stamp <= ${to} AND p_type ="F" LIMIT 5`;
    const results = await sequelize.query(query, { type: QueryTypes.SELECT });
    const body = results.shift();
    return new Response(JSON.stringify(body, null, 2));
}

async function esInTimeRange(esClient: Client): Promise<Response> {
    // magic random numbers
    const from = 1630000000 + Math.floor(Math.random() * 10000000);
    const to = from + Math.floor(Math.random() * 30000000);

    const body = {
        size: 5,
        query: {
            bool: {
                filter: [
                    { term: { p_type: "F" } },
                    {
                        range: {
                            time_stamp: {
                                gte: from,
                                lte: to,
                            },
                        },
                    },
                ],
            },
        },
    };
    const { body: results } = await esClient.search({ index: "prices", body });
    return new Response(JSON.stringify(results.hits.hits, null, 2));
}

export { now, es, randomByPrice, randomInTimeRange, esInTimeRange };
