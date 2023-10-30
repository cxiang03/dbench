import { Sequelize, QueryTypes } from "sequelize";

async function now(sequelize: Sequelize): Promise<Response> {
    const query = "SELECT NOW()";
    const results = await sequelize.query(query, { type: QueryTypes.SELECT });
    const body = results.shift();
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

export { now, randomByPrice, randomInTimeRange };
